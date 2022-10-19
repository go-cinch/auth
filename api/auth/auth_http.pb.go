// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.5.0
// - protoc             v3.17.3
// source: auth-proto/auth.proto

package auth

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationAuthCaptcha = "/auth.v1.Auth/Captcha"
const OperationAuthCreateAction = "/auth.v1.Auth/CreateAction"
const OperationAuthCreateRole = "/auth.v1.Auth/CreateRole"
const OperationAuthCreateUserGroup = "/auth.v1.Auth/CreateUserGroup"
const OperationAuthDeleteAction = "/auth.v1.Auth/DeleteAction"
const OperationAuthDeleteRole = "/auth.v1.Auth/DeleteRole"
const OperationAuthDeleteUser = "/auth.v1.Auth/DeleteUser"
const OperationAuthDeleteUserGroup = "/auth.v1.Auth/DeleteUserGroup"
const OperationAuthFindAction = "/auth.v1.Auth/FindAction"
const OperationAuthFindRole = "/auth.v1.Auth/FindRole"
const OperationAuthFindUser = "/auth.v1.Auth/FindUser"
const OperationAuthFindUserGroup = "/auth.v1.Auth/FindUserGroup"
const OperationAuthIdempotent = "/auth.v1.Auth/Idempotent"
const OperationAuthInfo = "/auth.v1.Auth/Info"
const OperationAuthLogin = "/auth.v1.Auth/Login"
const OperationAuthLogout = "/auth.v1.Auth/Logout"
const OperationAuthPermission = "/auth.v1.Auth/Permission"
const OperationAuthPwd = "/auth.v1.Auth/Pwd"
const OperationAuthRefresh = "/auth.v1.Auth/Refresh"
const OperationAuthRegister = "/auth.v1.Auth/Register"
const OperationAuthStatus = "/auth.v1.Auth/Status"
const OperationAuthUpdateAction = "/auth.v1.Auth/UpdateAction"
const OperationAuthUpdateRole = "/auth.v1.Auth/UpdateRole"
const OperationAuthUpdateUser = "/auth.v1.Auth/UpdateUser"
const OperationAuthUpdateUserGroup = "/auth.v1.Auth/UpdateUserGroup"

type AuthHTTPServer interface {
	Captcha(context.Context, *emptypb.Empty) (*CaptchaReply, error)
	CreateAction(context.Context, *CreateActionRequest) (*emptypb.Empty, error)
	CreateRole(context.Context, *CreateRoleRequest) (*emptypb.Empty, error)
	CreateUserGroup(context.Context, *CreateUserGroupRequest) (*emptypb.Empty, error)
	DeleteAction(context.Context, *IdsRequest) (*emptypb.Empty, error)
	DeleteRole(context.Context, *IdsRequest) (*emptypb.Empty, error)
	DeleteUser(context.Context, *IdsRequest) (*emptypb.Empty, error)
	DeleteUserGroup(context.Context, *IdsRequest) (*emptypb.Empty, error)
	FindAction(context.Context, *FindActionRequest) (*FindActionReply, error)
	FindRole(context.Context, *FindRoleRequest) (*FindRoleReply, error)
	FindUser(context.Context, *FindUserRequest) (*FindUserReply, error)
	FindUserGroup(context.Context, *FindUserGroupRequest) (*FindUserGroupReply, error)
	Idempotent(context.Context, *emptypb.Empty) (*IdempotentReply, error)
	Info(context.Context, *emptypb.Empty) (*InfoReply, error)
	Login(context.Context, *LoginRequest) (*LoginReply, error)
	Logout(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	Permission(context.Context, *PermissionRequest) (*PermissionReply, error)
	Pwd(context.Context, *PwdRequest) (*emptypb.Empty, error)
	Refresh(context.Context, *RefreshRequest) (*LoginReply, error)
	Register(context.Context, *RegisterRequest) (*emptypb.Empty, error)
	Status(context.Context, *StatusRequest) (*StatusReply, error)
	UpdateAction(context.Context, *UpdateActionRequest) (*emptypb.Empty, error)
	UpdateRole(context.Context, *UpdateRoleRequest) (*emptypb.Empty, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*emptypb.Empty, error)
	UpdateUserGroup(context.Context, *UpdateUserGroupRequest) (*emptypb.Empty, error)
}

func RegisterAuthHTTPServer(s *http.Server, srv AuthHTTPServer) {
	r := s.Route("/")
	r.POST("/register", _Auth_Register0_HTTP_Handler(srv))
	r.POST("/pwd", _Auth_Pwd0_HTTP_Handler(srv))
	r.POST("/login", _Auth_Login0_HTTP_Handler(srv))
	r.GET("/status", _Auth_Status0_HTTP_Handler(srv))
	r.GET("/captcha", _Auth_Captcha0_HTTP_Handler(srv))
	r.POST("/refresh", _Auth_Refresh0_HTTP_Handler(srv))
	r.POST("/logout", _Auth_Logout0_HTTP_Handler(srv))
	r.GET("/info", _Auth_Info0_HTTP_Handler(srv))
	r.GET("/idempotent", _Auth_Idempotent0_HTTP_Handler(srv))
	r.GET("/user", _Auth_FindUser0_HTTP_Handler(srv))
	r.PATCH("/user/{id}", _Auth_UpdateUser0_HTTP_Handler(srv))
	r.PUT("/user/{id}", _Auth_UpdateUser1_HTTP_Handler(srv))
	r.DELETE("/user/{ids}", _Auth_DeleteUser0_HTTP_Handler(srv))
	r.POST("/permission", _Auth_Permission0_HTTP_Handler(srv))
	r.POST("/action", _Auth_CreateAction0_HTTP_Handler(srv))
	r.GET("/action", _Auth_FindAction0_HTTP_Handler(srv))
	r.PATCH("/action/{id}", _Auth_UpdateAction0_HTTP_Handler(srv))
	r.PUT("/action/{id}", _Auth_UpdateAction1_HTTP_Handler(srv))
	r.DELETE("/action/{ids}", _Auth_DeleteAction0_HTTP_Handler(srv))
	r.POST("/role", _Auth_CreateRole0_HTTP_Handler(srv))
	r.GET("/role", _Auth_FindRole0_HTTP_Handler(srv))
	r.PATCH("/role/{id}", _Auth_UpdateRole0_HTTP_Handler(srv))
	r.PUT("/role/{id}", _Auth_UpdateRole1_HTTP_Handler(srv))
	r.DELETE("/role/{ids}", _Auth_DeleteRole0_HTTP_Handler(srv))
	r.POST("/user/group", _Auth_CreateUserGroup0_HTTP_Handler(srv))
	r.GET("/user/group", _Auth_FindUserGroup0_HTTP_Handler(srv))
	r.PATCH("/user/group/{id}", _Auth_UpdateUserGroup0_HTTP_Handler(srv))
	r.PUT("/user/group/{id}", _Auth_UpdateUserGroup1_HTTP_Handler(srv))
	r.DELETE("/user/group/{ids}", _Auth_DeleteUserGroup0_HTTP_Handler(srv))
}

func _Auth_Register0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in RegisterRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthRegister)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Register(ctx, req.(*RegisterRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_Pwd0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in PwdRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthPwd)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Pwd(ctx, req.(*PwdRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_Login0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in LoginRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Login(ctx, req.(*LoginRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LoginReply)
		return ctx.Result(200, reply)
	}
}

func _Auth_Status0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in StatusRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthStatus)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Status(ctx, req.(*StatusRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*StatusReply)
		return ctx.Result(200, reply)
	}
}

func _Auth_Captcha0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthCaptcha)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Captcha(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CaptchaReply)
		return ctx.Result(200, reply)
	}
}

func _Auth_Refresh0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in RefreshRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthRefresh)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Refresh(ctx, req.(*RefreshRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LoginReply)
		return ctx.Result(200, reply)
	}
}

func _Auth_Logout0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthLogout)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Logout(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_Info0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthInfo)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Info(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*InfoReply)
		return ctx.Result(200, reply)
	}
}

func _Auth_Idempotent0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthIdempotent)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Idempotent(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*IdempotentReply)
		return ctx.Result(200, reply)
	}
}

func _Auth_FindUser0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in FindUserRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthFindUser)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.FindUser(ctx, req.(*FindUserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*FindUserReply)
		return ctx.Result(200, reply)
	}
}

func _Auth_UpdateUser0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateUserRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthUpdateUser)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateUser(ctx, req.(*UpdateUserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_UpdateUser1_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateUserRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthUpdateUser)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateUser(ctx, req.(*UpdateUserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_DeleteUser0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in IdsRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthDeleteUser)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteUser(ctx, req.(*IdsRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_Permission0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in PermissionRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthPermission)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Permission(ctx, req.(*PermissionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*PermissionReply)
		return ctx.Result(200, reply)
	}
}

func _Auth_CreateAction0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateActionRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthCreateAction)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateAction(ctx, req.(*CreateActionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_FindAction0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in FindActionRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthFindAction)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.FindAction(ctx, req.(*FindActionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*FindActionReply)
		return ctx.Result(200, reply)
	}
}

func _Auth_UpdateAction0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateActionRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthUpdateAction)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateAction(ctx, req.(*UpdateActionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_UpdateAction1_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateActionRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthUpdateAction)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateAction(ctx, req.(*UpdateActionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_DeleteAction0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in IdsRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthDeleteAction)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteAction(ctx, req.(*IdsRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_CreateRole0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateRoleRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthCreateRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateRole(ctx, req.(*CreateRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_FindRole0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in FindRoleRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthFindRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.FindRole(ctx, req.(*FindRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*FindRoleReply)
		return ctx.Result(200, reply)
	}
}

func _Auth_UpdateRole0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateRoleRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthUpdateRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateRole(ctx, req.(*UpdateRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_UpdateRole1_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateRoleRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthUpdateRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateRole(ctx, req.(*UpdateRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_DeleteRole0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in IdsRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthDeleteRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteRole(ctx, req.(*IdsRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_CreateUserGroup0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateUserGroupRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthCreateUserGroup)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateUserGroup(ctx, req.(*CreateUserGroupRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_FindUserGroup0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in FindUserGroupRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthFindUserGroup)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.FindUserGroup(ctx, req.(*FindUserGroupRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*FindUserGroupReply)
		return ctx.Result(200, reply)
	}
}

func _Auth_UpdateUserGroup0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateUserGroupRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthUpdateUserGroup)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateUserGroup(ctx, req.(*UpdateUserGroupRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_UpdateUserGroup1_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateUserGroupRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthUpdateUserGroup)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateUserGroup(ctx, req.(*UpdateUserGroupRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Auth_DeleteUserGroup0_HTTP_Handler(srv AuthHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in IdsRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthDeleteUserGroup)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteUserGroup(ctx, req.(*IdsRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

type AuthHTTPClient interface {
	Captcha(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *CaptchaReply, err error)
	CreateAction(ctx context.Context, req *CreateActionRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	CreateRole(ctx context.Context, req *CreateRoleRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	CreateUserGroup(ctx context.Context, req *CreateUserGroupRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	DeleteAction(ctx context.Context, req *IdsRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	DeleteRole(ctx context.Context, req *IdsRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	DeleteUser(ctx context.Context, req *IdsRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	DeleteUserGroup(ctx context.Context, req *IdsRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	FindAction(ctx context.Context, req *FindActionRequest, opts ...http.CallOption) (rsp *FindActionReply, err error)
	FindRole(ctx context.Context, req *FindRoleRequest, opts ...http.CallOption) (rsp *FindRoleReply, err error)
	FindUser(ctx context.Context, req *FindUserRequest, opts ...http.CallOption) (rsp *FindUserReply, err error)
	FindUserGroup(ctx context.Context, req *FindUserGroupRequest, opts ...http.CallOption) (rsp *FindUserGroupReply, err error)
	Idempotent(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *IdempotentReply, err error)
	Info(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *InfoReply, err error)
	Login(ctx context.Context, req *LoginRequest, opts ...http.CallOption) (rsp *LoginReply, err error)
	Logout(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	Permission(ctx context.Context, req *PermissionRequest, opts ...http.CallOption) (rsp *PermissionReply, err error)
	Pwd(ctx context.Context, req *PwdRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	Refresh(ctx context.Context, req *RefreshRequest, opts ...http.CallOption) (rsp *LoginReply, err error)
	Register(ctx context.Context, req *RegisterRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	Status(ctx context.Context, req *StatusRequest, opts ...http.CallOption) (rsp *StatusReply, err error)
	UpdateAction(ctx context.Context, req *UpdateActionRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	UpdateRole(ctx context.Context, req *UpdateRoleRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	UpdateUser(ctx context.Context, req *UpdateUserRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	UpdateUserGroup(ctx context.Context, req *UpdateUserGroupRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
}

type AuthHTTPClientImpl struct {
	cc *http.Client
}

func NewAuthHTTPClient(client *http.Client) AuthHTTPClient {
	return &AuthHTTPClientImpl{client}
}

func (c *AuthHTTPClientImpl) Captcha(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*CaptchaReply, error) {
	var out CaptchaReply
	pattern := "/captcha"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAuthCaptcha))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) CreateAction(ctx context.Context, in *CreateActionRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/action"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthCreateAction))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) CreateRole(ctx context.Context, in *CreateRoleRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/role"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthCreateRole))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) CreateUserGroup(ctx context.Context, in *CreateUserGroupRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/user/group"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthCreateUserGroup))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) DeleteAction(ctx context.Context, in *IdsRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/action/{ids}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAuthDeleteAction))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) DeleteRole(ctx context.Context, in *IdsRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/role/{ids}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAuthDeleteRole))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) DeleteUser(ctx context.Context, in *IdsRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/user/{ids}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAuthDeleteUser))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) DeleteUserGroup(ctx context.Context, in *IdsRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/user/group/{ids}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAuthDeleteUserGroup))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) FindAction(ctx context.Context, in *FindActionRequest, opts ...http.CallOption) (*FindActionReply, error) {
	var out FindActionReply
	pattern := "/action"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAuthFindAction))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) FindRole(ctx context.Context, in *FindRoleRequest, opts ...http.CallOption) (*FindRoleReply, error) {
	var out FindRoleReply
	pattern := "/role"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAuthFindRole))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) FindUser(ctx context.Context, in *FindUserRequest, opts ...http.CallOption) (*FindUserReply, error) {
	var out FindUserReply
	pattern := "/user"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAuthFindUser))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) FindUserGroup(ctx context.Context, in *FindUserGroupRequest, opts ...http.CallOption) (*FindUserGroupReply, error) {
	var out FindUserGroupReply
	pattern := "/user/group"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAuthFindUserGroup))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) Idempotent(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*IdempotentReply, error) {
	var out IdempotentReply
	pattern := "/idempotent"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAuthIdempotent))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) Info(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*InfoReply, error) {
	var out InfoReply
	pattern := "/info"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAuthInfo))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) Login(ctx context.Context, in *LoginRequest, opts ...http.CallOption) (*LoginReply, error) {
	var out LoginReply
	pattern := "/login"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) Logout(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/logout"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthLogout))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) Permission(ctx context.Context, in *PermissionRequest, opts ...http.CallOption) (*PermissionReply, error) {
	var out PermissionReply
	pattern := "/permission"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthPermission))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) Pwd(ctx context.Context, in *PwdRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/pwd"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthPwd))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) Refresh(ctx context.Context, in *RefreshRequest, opts ...http.CallOption) (*LoginReply, error) {
	var out LoginReply
	pattern := "/refresh"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthRefresh))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) Register(ctx context.Context, in *RegisterRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/register"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthRegister))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) Status(ctx context.Context, in *StatusRequest, opts ...http.CallOption) (*StatusReply, error) {
	var out StatusReply
	pattern := "/status"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAuthStatus))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) UpdateAction(ctx context.Context, in *UpdateActionRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/action/{id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthUpdateAction))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) UpdateRole(ctx context.Context, in *UpdateRoleRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/role/{id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthUpdateRole))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/user/{id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthUpdateUser))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *AuthHTTPClientImpl) UpdateUserGroup(ctx context.Context, in *UpdateUserGroupRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/user/group/{id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthUpdateUserGroup))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}