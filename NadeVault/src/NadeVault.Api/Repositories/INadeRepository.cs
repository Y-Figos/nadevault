namespace NadeVault.Api.Repo;

using NadeVault.Api.Models;
public interface INadeRepository
{
    public List<Nade> GetAll();
}