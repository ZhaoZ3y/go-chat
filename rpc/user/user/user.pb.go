// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.19.4
// source: user.proto

package user

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 用户信息
type User struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Username      string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Email         string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Nickname      string                 `protobuf:"bytes,4,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Avatar        string                 `protobuf:"bytes,5,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Phone         string                 `protobuf:"bytes,6,opt,name=phone,proto3" json:"phone,omitempty"`
	Status        int32                  `protobuf:"varint,7,opt,name=status,proto3" json:"status,omitempty"` // 1:正常 2:禁用
	CreateAt      int64                  `protobuf:"varint,8,opt,name=create_at,json=createAt,proto3" json:"create_at,omitempty"`
	UpdateAt      int64                  `protobuf:"varint,9,opt,name=update_at,json=updateAt,proto3" json:"update_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_user_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *User) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *User) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *User) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *User) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *User) GetCreateAt() int64 {
	if x != nil {
		return x.CreateAt
	}
	return 0
}

func (x *User) GetUpdateAt() int64 {
	if x != nil {
		return x.UpdateAt
	}
	return 0
}

// 用户注册请求
type RegisterRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Username      string                 `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Email         string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Nickname      string                 `protobuf:"bytes,4,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Avatar        string                 `protobuf:"bytes,5,opt,name=avatar,proto3" json:"avatar,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	mi := &file_user_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *RegisterRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *RegisterRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *RegisterRequest) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *RegisterRequest) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

// 用户注册响应
type RegisterResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterResponse) Reset() {
	*x = RegisterResponse{}
	mi := &file_user_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterResponse) ProtoMessage() {}

func (x *RegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterResponse.ProtoReflect.Descriptor instead.
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// 用户登录请求
type LoginRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Username      string                 `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	mi := &file_user_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{3}
}

func (x *LoginRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// 用户登录响应
type LoginResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AccessToken   string                 `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	RefreshToken  string                 `protobuf:"bytes,3,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	mi := &file_user_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{4}
}

func (x *LoginResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *LoginResponse) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

// 获取用户信息请求
type GetUserInfoRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserInfoRequest) Reset() {
	*x = GetUserInfoRequest{}
	mi := &file_user_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserInfoRequest) ProtoMessage() {}

func (x *GetUserInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserInfoRequest.ProtoReflect.Descriptor instead.
func (*GetUserInfoRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{5}
}

func (x *GetUserInfoRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

// 获取用户信息响应
type GetUserInfoResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserInfo      *User                  `protobuf:"bytes,1,opt,name=user_info,json=userInfo,proto3" json:"user_info,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserInfoResponse) Reset() {
	*x = GetUserInfoResponse{}
	mi := &file_user_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserInfoResponse) ProtoMessage() {}

func (x *GetUserInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserInfoResponse.ProtoReflect.Descriptor instead.
func (*GetUserInfoResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{6}
}

func (x *GetUserInfoResponse) GetUserInfo() *User {
	if x != nil {
		return x.UserInfo
	}
	return nil
}

// 更新用户信息请求
type UpdateUserInfoRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Nickname      string                 `protobuf:"bytes,2,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Avatar        string                 `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Phone         string                 `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty"`
	Email         string                 `protobuf:"bytes,5,opt,name=email,proto3" json:"email,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateUserInfoRequest) Reset() {
	*x = UpdateUserInfoRequest{}
	mi := &file_user_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateUserInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserInfoRequest) ProtoMessage() {}

func (x *UpdateUserInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserInfoRequest.ProtoReflect.Descriptor instead.
func (*UpdateUserInfoRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateUserInfoRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UpdateUserInfoRequest) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *UpdateUserInfoRequest) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *UpdateUserInfoRequest) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *UpdateUserInfoRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

// 更新用户信息响应
type UpdateUserInfoResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateUserInfoResponse) Reset() {
	*x = UpdateUserInfoResponse{}
	mi := &file_user_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateUserInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserInfoResponse) ProtoMessage() {}

func (x *UpdateUserInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserInfoResponse.ProtoReflect.Descriptor instead.
func (*UpdateUserInfoResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{8}
}

func (x *UpdateUserInfoResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// 搜索用户请求
type SearchUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Keyword       string                 `protobuf:"bytes,1,opt,name=keyword,proto3" json:"keyword,omitempty"`
	CurrentUserId int64                  `protobuf:"varint,2,opt,name=currentUserId,proto3" json:"currentUserId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchUserRequest) Reset() {
	*x = SearchUserRequest{}
	mi := &file_user_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchUserRequest) ProtoMessage() {}

func (x *SearchUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchUserRequest.ProtoReflect.Descriptor instead.
func (*SearchUserRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{9}
}

func (x *SearchUserRequest) GetKeyword() string {
	if x != nil {
		return x.Keyword
	}
	return ""
}

func (x *SearchUserRequest) GetCurrentUserId() int64 {
	if x != nil {
		return x.CurrentUserId
	}
	return 0
}

// 搜索用户响应
type SearchUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Users         []*User                `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
	Total         int64                  `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchUserResponse) Reset() {
	*x = SearchUserResponse{}
	mi := &file_user_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchUserResponse) ProtoMessage() {}

func (x *SearchUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchUserResponse.ProtoReflect.Descriptor instead.
func (*SearchUserResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{10}
}

func (x *SearchUserResponse) GetUsers() []*User {
	if x != nil {
		return x.Users
	}
	return nil
}

func (x *SearchUserResponse) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

// 刷新Token请求
type RefreshTokenRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RefreshToken  string                 `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RefreshTokenRequest) Reset() {
	*x = RefreshTokenRequest{}
	mi := &file_user_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RefreshTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshTokenRequest) ProtoMessage() {}

func (x *RefreshTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshTokenRequest.ProtoReflect.Descriptor instead.
func (*RefreshTokenRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{11}
}

func (x *RefreshTokenRequest) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

// 刷新Token响应
type RefreshTokenResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AccessToken   string                 `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	RefreshToken  string                 `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RefreshTokenResponse) Reset() {
	*x = RefreshTokenResponse{}
	mi := &file_user_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RefreshTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshTokenResponse) ProtoMessage() {}

func (x *RefreshTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshTokenResponse.ProtoReflect.Descriptor instead.
func (*RefreshTokenResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{12}
}

func (x *RefreshTokenResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *RefreshTokenResponse) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

// 修改密码请求
type ChangePasswordRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	OldPassword   string                 `protobuf:"bytes,2,opt,name=old_password,json=oldPassword,proto3" json:"old_password,omitempty"`
	NewPassword   string                 `protobuf:"bytes,3,opt,name=new_password,json=newPassword,proto3" json:"new_password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChangePasswordRequest) Reset() {
	*x = ChangePasswordRequest{}
	mi := &file_user_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChangePasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangePasswordRequest) ProtoMessage() {}

func (x *ChangePasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangePasswordRequest.ProtoReflect.Descriptor instead.
func (*ChangePasswordRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{13}
}

func (x *ChangePasswordRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ChangePasswordRequest) GetOldPassword() string {
	if x != nil {
		return x.OldPassword
	}
	return ""
}

func (x *ChangePasswordRequest) GetNewPassword() string {
	if x != nil {
		return x.NewPassword
	}
	return ""
}

// 修改密码响应
type ChangePasswordResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChangePasswordResponse) Reset() {
	*x = ChangePasswordResponse{}
	mi := &file_user_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChangePasswordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangePasswordResponse) ProtoMessage() {}

func (x *ChangePasswordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangePasswordResponse.ProtoReflect.Descriptor instead.
func (*ChangePasswordResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{14}
}

func (x *ChangePasswordResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_user_proto protoreflect.FileDescriptor

const file_user_proto_rawDesc = "" +
	"\n" +
	"\n" +
	"user.proto\x12\x04user\"\xe4\x01\n" +
	"\x04User\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x1a\n" +
	"\busername\x18\x02 \x01(\tR\busername\x12\x14\n" +
	"\x05email\x18\x03 \x01(\tR\x05email\x12\x1a\n" +
	"\bnickname\x18\x04 \x01(\tR\bnickname\x12\x16\n" +
	"\x06avatar\x18\x05 \x01(\tR\x06avatar\x12\x14\n" +
	"\x05phone\x18\x06 \x01(\tR\x05phone\x12\x16\n" +
	"\x06status\x18\a \x01(\x05R\x06status\x12\x1b\n" +
	"\tcreate_at\x18\b \x01(\x03R\bcreateAt\x12\x1b\n" +
	"\tupdate_at\x18\t \x01(\x03R\bupdateAt\"\x93\x01\n" +
	"\x0fRegisterRequest\x12\x1a\n" +
	"\busername\x18\x01 \x01(\tR\busername\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\x12\x14\n" +
	"\x05email\x18\x03 \x01(\tR\x05email\x12\x1a\n" +
	"\bnickname\x18\x04 \x01(\tR\bnickname\x12\x16\n" +
	"\x06avatar\x18\x05 \x01(\tR\x06avatar\",\n" +
	"\x10RegisterResponse\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\"F\n" +
	"\fLoginRequest\x12\x1a\n" +
	"\busername\x18\x01 \x01(\tR\busername\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\"W\n" +
	"\rLoginResponse\x12!\n" +
	"\faccess_token\x18\x02 \x01(\tR\vaccessToken\x12#\n" +
	"\rrefresh_token\x18\x03 \x01(\tR\frefreshToken\"-\n" +
	"\x12GetUserInfoRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x03R\x06userId\">\n" +
	"\x13GetUserInfoResponse\x12'\n" +
	"\tuser_info\x18\x01 \x01(\v2\n" +
	".user.UserR\buserInfo\"\x90\x01\n" +
	"\x15UpdateUserInfoRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x03R\x06userId\x12\x1a\n" +
	"\bnickname\x18\x02 \x01(\tR\bnickname\x12\x16\n" +
	"\x06avatar\x18\x03 \x01(\tR\x06avatar\x12\x14\n" +
	"\x05phone\x18\x04 \x01(\tR\x05phone\x12\x14\n" +
	"\x05email\x18\x05 \x01(\tR\x05email\"2\n" +
	"\x16UpdateUserInfoResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"S\n" +
	"\x11SearchUserRequest\x12\x18\n" +
	"\akeyword\x18\x01 \x01(\tR\akeyword\x12$\n" +
	"\rcurrentUserId\x18\x02 \x01(\x03R\rcurrentUserId\"L\n" +
	"\x12SearchUserResponse\x12 \n" +
	"\x05users\x18\x01 \x03(\v2\n" +
	".user.UserR\x05users\x12\x14\n" +
	"\x05total\x18\x02 \x01(\x03R\x05total\":\n" +
	"\x13RefreshTokenRequest\x12#\n" +
	"\rrefresh_token\x18\x01 \x01(\tR\frefreshToken\"^\n" +
	"\x14RefreshTokenResponse\x12!\n" +
	"\faccess_token\x18\x01 \x01(\tR\vaccessToken\x12#\n" +
	"\rrefresh_token\x18\x02 \x01(\tR\frefreshToken\"v\n" +
	"\x15ChangePasswordRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x03R\x06userId\x12!\n" +
	"\fold_password\x18\x02 \x01(\tR\voldPassword\x12!\n" +
	"\fnew_password\x18\x03 \x01(\tR\vnewPassword\"2\n" +
	"\x16ChangePasswordResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage2\xe0\x03\n" +
	"\vUserService\x129\n" +
	"\bRegister\x12\x15.user.RegisterRequest\x1a\x16.user.RegisterResponse\x120\n" +
	"\x05Login\x12\x12.user.LoginRequest\x1a\x13.user.LoginResponse\x12B\n" +
	"\vGetUserInfo\x12\x18.user.GetUserInfoRequest\x1a\x19.user.GetUserInfoResponse\x12K\n" +
	"\x0eUpdateUserInfo\x12\x1b.user.UpdateUserInfoRequest\x1a\x1c.user.UpdateUserInfoResponse\x12?\n" +
	"\n" +
	"SearchUser\x12\x17.user.SearchUserRequest\x1a\x18.user.SearchUserResponse\x12E\n" +
	"\fRefreshToken\x12\x19.user.RefreshTokenRequest\x1a\x1a.user.RefreshTokenResponse\x12K\n" +
	"\x0eChangePassword\x12\x1b.user.ChangePasswordRequest\x1a\x1c.user.ChangePasswordResponseB\bZ\x06./userb\x06proto3"

var (
	file_user_proto_rawDescOnce sync.Once
	file_user_proto_rawDescData []byte
)

func file_user_proto_rawDescGZIP() []byte {
	file_user_proto_rawDescOnce.Do(func() {
		file_user_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_user_proto_rawDesc), len(file_user_proto_rawDesc)))
	})
	return file_user_proto_rawDescData
}

var file_user_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_user_proto_goTypes = []any{
	(*User)(nil),                   // 0: user.User
	(*RegisterRequest)(nil),        // 1: user.RegisterRequest
	(*RegisterResponse)(nil),       // 2: user.RegisterResponse
	(*LoginRequest)(nil),           // 3: user.LoginRequest
	(*LoginResponse)(nil),          // 4: user.LoginResponse
	(*GetUserInfoRequest)(nil),     // 5: user.GetUserInfoRequest
	(*GetUserInfoResponse)(nil),    // 6: user.GetUserInfoResponse
	(*UpdateUserInfoRequest)(nil),  // 7: user.UpdateUserInfoRequest
	(*UpdateUserInfoResponse)(nil), // 8: user.UpdateUserInfoResponse
	(*SearchUserRequest)(nil),      // 9: user.SearchUserRequest
	(*SearchUserResponse)(nil),     // 10: user.SearchUserResponse
	(*RefreshTokenRequest)(nil),    // 11: user.RefreshTokenRequest
	(*RefreshTokenResponse)(nil),   // 12: user.RefreshTokenResponse
	(*ChangePasswordRequest)(nil),  // 13: user.ChangePasswordRequest
	(*ChangePasswordResponse)(nil), // 14: user.ChangePasswordResponse
}
var file_user_proto_depIdxs = []int32{
	0,  // 0: user.GetUserInfoResponse.user_info:type_name -> user.User
	0,  // 1: user.SearchUserResponse.users:type_name -> user.User
	1,  // 2: user.UserService.Register:input_type -> user.RegisterRequest
	3,  // 3: user.UserService.Login:input_type -> user.LoginRequest
	5,  // 4: user.UserService.GetUserInfo:input_type -> user.GetUserInfoRequest
	7,  // 5: user.UserService.UpdateUserInfo:input_type -> user.UpdateUserInfoRequest
	9,  // 6: user.UserService.SearchUser:input_type -> user.SearchUserRequest
	11, // 7: user.UserService.RefreshToken:input_type -> user.RefreshTokenRequest
	13, // 8: user.UserService.ChangePassword:input_type -> user.ChangePasswordRequest
	2,  // 9: user.UserService.Register:output_type -> user.RegisterResponse
	4,  // 10: user.UserService.Login:output_type -> user.LoginResponse
	6,  // 11: user.UserService.GetUserInfo:output_type -> user.GetUserInfoResponse
	8,  // 12: user.UserService.UpdateUserInfo:output_type -> user.UpdateUserInfoResponse
	10, // 13: user.UserService.SearchUser:output_type -> user.SearchUserResponse
	12, // 14: user.UserService.RefreshToken:output_type -> user.RefreshTokenResponse
	14, // 15: user.UserService.ChangePassword:output_type -> user.ChangePasswordResponse
	9,  // [9:16] is the sub-list for method output_type
	2,  // [2:9] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_user_proto_init() }
func file_user_proto_init() {
	if File_user_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_user_proto_rawDesc), len(file_user_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_proto_goTypes,
		DependencyIndexes: file_user_proto_depIdxs,
		MessageInfos:      file_user_proto_msgTypes,
	}.Build()
	File_user_proto = out.File
	file_user_proto_goTypes = nil
	file_user_proto_depIdxs = nil
}
