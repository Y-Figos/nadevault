namespace NadeVault.Api.Models;
public enum CsMap
{
    Mirage,
    Inferno,
    Ancient,
    Anubis,
    Dust2,
    Overpass,
    Cache,
    Nuke,
    Train,
    Cobblestone

}

public enum NadeType
{
    Smoke,
    Moly,
    Flash,
    Frag,
}

public enum MouseClick
{
    Mouse1,
    Mouse2,
    Both,
}

public class Nade
{
    public int Id {get;set;}
    public string Name {get;set;} = string.Empty;
    public CsMap Map {get;set;}
    public NadeType Type {get;set;} //SMOKE,MOLY,FLASH,FRAG
    //Input Modifiers
    public MouseClick MouseClick {get;set;} // MOUSE1 | MOUSE2 | BOTH
    public bool IsJumping {get;set;}
    public bool IsRunning {get;set;}
    public bool IsWalking {get;set;}
    //Imagens
    public List<string>? AimAtUrl {get;set;}
    public List<string>? PositionUrl {get;set;}
    public List<string>? TargetUrl {get;set;}
}

