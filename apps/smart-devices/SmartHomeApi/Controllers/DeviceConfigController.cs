using Microsoft.AspNetCore.Mvc;
using SmartHomeApi.Models;
using System.Collections.Generic;

namespace SmartHomeApi.Controllers
{
    [ApiController]
    [Route("api/v1/devices/{id}/config")]
    public class DeviceConfigController : ControllerBase
    {
        private static readonly Dictionary<string, DeviceConfig> Configs = new();

        [HttpGet]
        public IActionResult GetConfig(string id)
        {
            if (Configs.TryGetValue(id, out var config))
                return Ok(config);

            return NotFound(new { error = "Device not found", code = "DEVICE_NOT_FOUND" });
        }

        [HttpPut]
        public IActionResult UpdateConfig(string id, [FromBody] DeviceConfig updated)
        {
            updated.DeviceId = id;
            updated.LastUpdated = System.DateTime.UtcNow;
            Configs[id] = updated;
            return Ok(updated);
        }
    }
}