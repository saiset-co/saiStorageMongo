package api

func InitAPI() {
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
