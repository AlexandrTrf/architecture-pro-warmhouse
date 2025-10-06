using Microsoft.AspNetCore.Mvc;

namespace SmartHomeApi.Controllers
{
    [ApiController]
    [Route("health")]
    public class HealthController : ControllerBase
    {
        [HttpGet]
        public IActionResult GetHealth()
        {
            return Ok(new { status = "ok", service = "device-management" });
        }
    }
}