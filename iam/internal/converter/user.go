package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/alexander-kartavtsev/starship/iam/internal/model"
	commonV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/common/v1"
)

func UserToProto(user *model.User) *commonV1.User {
	if user == nil {
		return nil
	}
	return &commonV1.User{
		Uuid:      user.Uuid,
		Info:      UserInfoToProto(user.Info),
		CreatedAt: timestamppb.New(*user.CreatedAt),
		UpdatedAt: timestamppb.New(*user.UpdatedAt),
	}
}

func UserInfoToProto(info *model.UserInfo) *commonV1.UserInfo {
	if info == nil {
		return nil
	}
	return &commonV1.UserInfo{
		Login:               info.Login,
		Email:               info.Email,
		NotificationMethods: NotificationMethodsToProto(info.NotificationMethods),
	}
}

func NotificationMethodsToProto(methods []*model.NotificationMethod) []*commonV1.NotificationMethod {
	if len(methods) == 0 {
		return []*commonV1.NotificationMethod{}
	}
	var m []*commonV1.NotificationMethod
	for _, methodsItem := range methods {
		mItem := commonV1.NotificationMethod{
			ProviderName: methodsItem.ProviderName,
			Target:       methodsItem.Target,
		}
		m = append(m, &mItem)
	}
	return m
}

func UserInfoToModel(info *commonV1.UserInfo) *model.UserInfo {
	if info == nil {
		return nil
	}
	return &model.UserInfo{
		Login:               info.Login,
		Email:               info.Email,
		NotificationMethods: NotificationMethodsToModel(info.NotificationMethods),
	}
}

func NotificationMethodsToModel(methods []*commonV1.NotificationMethod) []*model.NotificationMethod {
	if len(methods) == 0 {
		return []*model.NotificationMethod{}
	}
	var m []*model.NotificationMethod
	for _, methodsItem := range methods {
		mItem := model.NotificationMethod{
			ProviderName: methodsItem.ProviderName,
			Target:       methodsItem.Target,
		}
		m = append(m, &mItem)
	}
	return m
}
