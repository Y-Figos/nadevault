using Microsoft.AspNetCore.Mvc;
using NadeVault.Api.Models;
using NadeVault.Api.Repo;

namespace NadeVault.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class NadesController : ControllerBase
{

    private readonly INadeRepository _repository;

    public NadesController(INadeRepository repository)
    {
        _repository = repository;
    }

    [HttpGet]
    public IActionResult GetAll()
    {
        var nades = _repository.GetAll();

        return Ok(nades);
    }
}
