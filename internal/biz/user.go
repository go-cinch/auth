package biz

import (
	"auth/internal/conf"
	"auth/internal/pkg/jwt"
	"context"
	"fmt"
	"github.com/go-cinch/common/captcha"
	"github.com/go-cinch/common/constant"
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/utils"
	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           uint64          `json:"id,string"`
	CreatedAt    carbon.DateTime `json:"createdAt,string"`
	UpdatedAt    carbon.DateTime `json:"updatedAt,string"`
	RoleId       uint64          `json:"roleId,string"`
	Role         Role            `json:"role"`
	Action       string          `json:"action"`
	Actions      []Action        `json:"actions"`
	Username     string          `json:"username"`
	Code         string          `json:"Code"`
	Password     string          `json:"password"`
	OldPassword  string          `json:"-"`
	NewPassword  string          `json:"-"`
	Mobile       string          `json:"mobile"`
	Avatar       string          `json:"avatar"`
	Nickname     string          `json:"nickname"`
	Introduction string          `json:"introduction"`
	LastLogin    carbon.DateTime `json:"lastLogin,string,omitempty"`
	Locked       uint64          `json:"locked"`
	LockExpire   int64           `json:"lockExpire"`
	LockMsg      string          `json:"lockMsg"`
	Wrong        int64           `json:"wrong"`
	Captcha      Captcha         `json:"-"`
}

type UserInfo struct {
	Id           uint64 `json:"id,string"`
	Username     string `json:"username"`
	Code         string `json:"code"`
	Mobile       string `json:"mobile"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
}

type FindUser struct {
	Page           page.Page `json:"page"`
	StartCreatedAt *string   `json:"startCreatedAt"`
	EndCreatedAt   *string   `json:"endCreatedAt"`
	StartUpdatedAt *string   `json:"startUpdatedAt"`
	EndUpdatedAt   *string   `json:"endUpdatedAt"`
	Username       *string   `json:"username"`
	Code           *string   `json:"code"`
	Mobile         *string   `json:"mobile"`
	Locked         *uint64   `json:"locked"`
}

type UpdateUser struct {
	Id           *uint64 `json:"id,string,omitempty"`
	Action       *string `json:"action,omitempty"`
	Username     *string `json:"username,omitempty"`
	Password     *string `json:"password,omitempty"`
	Mobile       *string `json:"mobile,omitempty"`
	Avatar       *string `json:"avatar,omitempty"`
	Nickname     *string `json:"nickname,omitempty"`
	Introduction *string `json:"introduction,omitempty"`
	Locked       *uint64 `json:"locked,omitempty"`
	LockExpire   *int64  `json:"lockExpire,omitempty"`
}

type Login struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	CaptchaId     string `json:"captchaId"`
	CaptchaAnswer string `json:"captchaAnswer"`
}

type LoginTime struct {
	Username  string          `json:"username"`
	LastLogin carbon.DateTime `json:"lastLogin"`
}

type LoginToken struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

type UserStatus struct {
	Id          uint64  `json:"id,string"`
	Code        string  `json:"code"`
	Password    string  `json:"password"`
	Wrong       int64   `json:"wrong"`
	Locked      uint64  `json:"locked"`
	LockExpire  int64   `json:"lockExpire"`
	NeedCaptcha bool    `json:"needCaptcha"`
	Captcha     Captcha `json:"captcha"`
}

type Captcha struct {
	Id  string `json:"id"`
	Img string `json:"img"`
}

type UserRepo interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	Find(ctx context.Context, condition *FindUser) []User
	Create(ctx context.Context, item *User) error
	Update(ctx context.Context, item *UpdateUser) error
	Delete(ctx context.Context, ids ...uint64) error
	LastLogin(ctx context.Context, id uint64) error
	WrongPwd(ctx context.Context, req LoginTime) error
	UpdatePassword(ctx context.Context, item *User) error
	IdExists(ctx context.Context, id uint64) error
	GetByCode(ctx context.Context, code string) (*User, error)
}

type UserUseCase struct {
	c     *conf.Bootstrap
	repo  UserRepo
	tx    Transaction
	cache Cache
}

func NewUserUseCase(c *conf.Bootstrap, repo UserRepo, tx Transaction, cache Cache) *UserUseCase {
	return &UserUseCase{c: c, repo: repo, tx: tx, cache: cache.WithPrefix("auth_user")}
}

func (uc *UserUseCase) Create(ctx context.Context, item *User) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			item.Password = genPwd(item.Password)
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *UserUseCase) Update(ctx context.Context, item *UpdateUser) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			if item.Password != nil {
				pwd := genPwd(*item.Password)
				item.Password = &pwd
			}
			err = uc.repo.Update(ctx, item)
			return
		})
	})
}

func (uc *UserUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			info, err := uc.InfoFromCtx(ctx)
			if err != nil {
				return
			}
			if funk.ContainsUInt64(ids, info.Id) {
				err = DeleteYourself
				return
			}
			err = uc.repo.Delete(ctx, ids...)
			return
		})
	})
}

func (uc *UserUseCase) Find(ctx context.Context, condition *FindUser) (rp []User) {
	rp = make([]User, 0)
	action := fmt.Sprintf("find_%s", utils.StructMd5(condition))
	str, ok, _, _ := uc.cache.Get(ctx, action, func(ctx context.Context) (string, bool) {
		return uc.find(ctx, action, condition)
	})
	if ok {
		utils.Json2Struct(&rp, str)
	}
	return
}

func (uc *UserUseCase) find(ctx context.Context, action string, condition *FindUser) (res string, ok bool) {
	// read data from db and write to cache
	rp := make([]User, 0)
	list := uc.repo.Find(ctx, condition)
	copierx.Copy(&rp, list)
	res = utils.Struct2Json(rp)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	ok = true
	return
}

func (uc *UserUseCase) InfoFromCtx(ctx context.Context) (rp *UserInfo, err error) {
	user := jwt.FromContext(ctx)
	return uc.Info(ctx, user.Code)
}

func (uc *UserUseCase) Info(ctx context.Context, code string) (rp *UserInfo, err error) {
	rp = &UserInfo{}
	action := fmt.Sprintf("info_%s", code)
	str, ok, lock, _ := uc.cache.Get(ctx, action, func(ctx context.Context) (string, bool) {
		return uc.info(ctx, action, code)
	})
	if ok {
		utils.Json2Struct(&rp, str)
	} else if !lock {
		err = TooManyRequests
		return
	}
	return
}

func (uc *UserUseCase) Login(ctx context.Context, item *Login) (rp *LoginToken, err error) {
	rp = &LoginToken{}
	err = uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			status, err := uc.Status(ctx, item.Username, false)
			if err != nil {
				return
			}
			if status.Id == constant.UI0 {
				err = NotFound("%s User.username: %s", RecordNotFound.Message, item.Username)
				return
			}
			// verify captcha
			if status.NeedCaptcha && !uc.VerifyCaptcha(ctx, item.CaptchaId, item.CaptchaAnswer) {
				err = InvalidCaptcha
				return
			}
			// user is locked
			if status.Locked == constant.UI1 {
				err = UserLocked
				return
			}
			// check password
			if ok := comparePwd(item.Password, status.Password); !ok {
				err = LoginFailed
				return
			}
			err = uc.repo.LastLogin(ctx, status.Id)
			if err != nil {
				return
			}
			authUser := jwt.User{
				Code: status.Code,
			}
			token, expireTime := authUser.CreateToken(uc.c.Auth.Jwt.Key, uc.c.Auth.Jwt.Expires)
			rp.Token = token
			rp.Expires = expireTime.ToDateTimeString()
			return
		})
	})
	return
}

func (uc *UserUseCase) WrongPwd(ctx context.Context, req LoginTime) (err error) {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.WrongPwd(ctx, req)
		})
	})
}

func (uc *UserUseCase) Pwd(ctx context.Context, item *User) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
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
			item.Password = genPwd(item.NewPassword)
			return uc.repo.UpdatePassword(ctx, item)
		})
	})
}

func (uc *UserUseCase) Status(ctx context.Context, username string, captcha bool) (rp *UserStatus, err error) {
	rp = &UserStatus{}
	action := fmt.Sprintf("status_%s", username)
	str, ok, lock, _ := uc.cache.Get(ctx, action, func(ctx context.Context) (string, bool) {
		return uc.status(ctx, action, username)
	})
	if ok {
		utils.Json2Struct(&rp, str)
		// TODO u can get max wrong count from env or dict
		if rp.Wrong >= constant.I3 {
			// need captcha
			rp.NeedCaptcha = true
			if captcha {
				rp.Captcha = uc.Captcha(ctx)
			}
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

func (uc *UserUseCase) VerifyCaptcha(ctx context.Context, id, answer string) bool {
	return captcha.New(
		captcha.WithRedis(uc.cache.Cache()),
		captcha.WithCtx(ctx),
	).Verify(id, answer)
}

func (uc *UserUseCase) status(ctx context.Context, action string, username string) (res string, ok bool) {
	// read data from db and write to cache
	rp := &UserStatus{}
	user, err := uc.repo.GetByUsername(ctx, username)
	if err != nil && !errors.Is(err, RecordNotFound) {
		return
	}
	copierx.Copy(&rp, user)
	res = utils.Struct2Json(rp)
	uc.cache.Set(ctx, action, res, errors.Is(err, RecordNotFound))
	ok = true
	return
}

func (uc *UserUseCase) info(ctx context.Context, action string, code string) (res string, ok bool) {
	// read data from db and write to cache
	rp := &UserInfo{}
	user, err := uc.repo.GetByCode(ctx, code)
	if err != nil && !errors.Is(err, RecordNotFound) {
		return
	}
	copierx.Copy(&rp, user)
	res = utils.Struct2Json(rp)
	uc.cache.Set(ctx, action, res, errors.Is(err, RecordNotFound))
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
