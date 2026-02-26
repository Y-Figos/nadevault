package nadevault

type CsMapProvider interface {
    GetCsMap() CsMap
}
func (r GetMapByIDRow) GetCsMap() CsMap { return r.CsMap }
func (r GetMapsRow) GetCsMap()    CsMap { return r.CsMap }

type NadeProvider interface {
	GetNadeParam() (Nade, CsMap)
}

func (r GetNadeByPublicIDRow) GetNadeParam() (Nade, CsMap) {return r.Nade, r.CsMap}
func (r GetNadeRow) GetNadeParam() (Nade, CsMap) {return r.Nade, r.CsMap}
func (r ListNadesByMapIDRow) GetNadeParam() (Nade, CsMap) {return r.Nade, r.CsMap}