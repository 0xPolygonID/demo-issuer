package models

type IdentityStatus string

type IdentityState struct {
	//StateID            int64
	//Identifier         string
	//State              *string
	//RootOfRoots        *string
	//ClaimsTreeRoot     *string
	//RevocationTreeRoot *string
}

//func (is IdentityState) ToTreeState() (circuits.TreeState, error) {
//	return BuildTreeState(is.State, is.ClaimsTreeRoot, is.RevocationTreeRoot,
//		is.RootOfRoots)
//}
