package appinfoUsecases

import "github.com/Kamila3820/go-shop-tutorial/modules/appinfo/appinfoRepositories"

type IAppinfoUsecase interface {
}

type appinfoUsecase struct {
	appinfoRepository appinfoRepositories.IAppinfoRepository
}

func AppinfoUsecase(appinfoRepository appinfoRepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{
		appinfoRepository: appinfoRepository,
	}
}
