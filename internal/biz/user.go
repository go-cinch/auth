package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-cinch/common/captcha"
	"github.com/go-cinch/common/constant"
	"github.com/golang-module/carbon/v2"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           uint64          `json:"id"`
	CreatedAt    carbon.DateTime `json:"createdAt"`
	UpdatedAt    carbon.DateTime `json:"updatedAt"`
	Username     string          `json:"username"`
	Password     string          `json:"password"`
	OldPassword  string          `json:"-"`
	NewPassword  string          `json:"-"`
	Mobile       string          `json:"mobile"`
	Avatar       string          `json:"avatar"`
	Nickname     string          `json:"nickname"`
	Introduction string          `json:"introduction"`
	Status       uint64          `json:"status"`
	LastLogin    carbon.DateTime `json:"lastLogin"`
	Locked       uint64          `json:"locked"`
	LockExpire   int64           `json:"lockExpire"`
	Wrong        int64           `json:"wrong"`
	Captcha      Captcha         `json:"-"`
}

type UserStatus struct {
	Wrong      int64   `json:"wrong"`
	Locked     uint64  `json:"locked"`
	LockExpire int64   `json:"lockExpire"`
	Captcha    Captcha `json:"captcha"`
}

type Captcha struct {
	Id  string `json:"id"`
	Img string `json:"img"`
}

type UserRepo interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, item *User) error
	UpdatePassword(ctx context.Context, item *User) error
}

type UserUseCase struct {
	repo  UserRepo
	tx    Transaction
	cache Cache
}

func NewUserUseCase(repo UserRepo, tx Transaction, cache Cache) *UserUseCase {
	cache.Register("auth_user_cache")
	return &UserUseCase{repo: repo, tx: tx, cache: cache}
}

func (uc *UserUseCase) Create(ctx context.Context, item *User) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		item.Password = genPwd(item.Password)
		return uc.repo.Create(ctx, item)
	})
}

func (uc *UserUseCase) UpdatePassword(ctx context.Context, item *User) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) (err error) {
		oldItem := &User{}
		oldItem, err = uc.repo.GetByUsername(ctx, item.Username)
		if err != nil {
			return
		}
		ok := comparePwd(item.OldPassword, oldItem.Password)
		if !ok {
			err = IncorrectPassword
			return
		}
		ok = comparePwd(item.NewPassword, oldItem.Password)
		if ok {
			err = SamePassword
			return
		}
		return uc.repo.UpdatePassword(ctx, item)
	})
}

func (uc *UserUseCase) Status(ctx context.Context, username string) (rp *UserStatus, err error) {
	rp = &UserStatus{}
	action := fmt.Sprintf("status_%s", username)
	str, ok, lock, _ := uc.cache.Get(ctx, action, func(ctx context.Context) (string, bool) {
		return uc.status(ctx, action, username)
	})
	if ok {
		json.Unmarshal([]byte(str), rp)
		// TODO u can get max wrong count from env or dict
		if rp.Wrong >= constant.I3 {
			// need captcha
			rp.Captcha = uc.Captcha(ctx)
		}
		timestamp := carbon.Now().Timestamp()
		if rp.Locked == constant.UI1 && rp.LockExpire > constant.I0 && timestamp >= rp.LockExpire {
			// unlock when lock time expiration
			rp.Locked = constant.UI0
		}
	} else if !lock {
		err = TooManyRequests
		return
	}
	return
}

func (uc *UserUseCase) Captcha(ctx context.Context) (rp Captcha) {
	rp.Id, rp.Img = captcha.New(
		captcha.WithRedis(uc.cache.Cache()),
		captcha.WithCtx(ctx),
	).Get()
	return
}

func (uc *UserUseCase) status(ctx context.Context, action string, username string) (res string, ok bool) {
	// read data from db and write to cache
	rp := &UserStatus{}
	user, err := uc.repo.GetByUsername(ctx, username)
	if err != nil && err != UserNotFound {
		return
	}
	copier.Copy(rp, user)
	bs, _ := json.Marshal(rp)
	res = string(bs)
	uc.cache.Set(ctx, action, res, err == UserNotFound)
	ok = true
	return
}

// generate password is irreversible due to the use of adaptive hash algorithm
func genPwd(str string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(hash)
}

// by comparing two string hashes, judge whether they are from the same plaintext
func comparePwd(str string, pwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(str)); err != nil {
		return false
	}
	return true
}
