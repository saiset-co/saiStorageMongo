package api

func InitAPI() {
	Ping()
	AddGetDataMethod()
	AddSaveDataMethod()
	AddUpdateDataMethod()
	AddUpsertDataMethod()
	AddRemoveDataMethod()
	Registration()
	Login()
	Logout()
	CreateRole()
}
