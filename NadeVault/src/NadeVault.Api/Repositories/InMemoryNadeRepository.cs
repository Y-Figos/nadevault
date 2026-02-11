namespace NadeVault.Api.Repo;

using NadeVault.Api.Models;

public class InMemoryNadeRepository : INadeRepository
{
    // Simulação de banco de dados
    private readonly List<Nade> _nades = new()
    {
        new Nade { Id = 1, Name = "Mirage Window", /*...*/ }
    };

    public List<Nade> GetAll()
    {
        return _nades;
    }
}