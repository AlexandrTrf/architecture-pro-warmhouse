using Microsoft.AspNetCore.Mvc;
using SmartHomeApi.Models;
using System.Collections.Generic;
using System.Linq;

namespace SmartHomeApi.Controllers
{
    [ApiController]
    [Route("api/v1/devices/{id}/status")]
    public class DeviceStatusController : ControllerBase
    {
        private static readonly Dictionary<string, DeviceStatus> Statuses = new();

        [HttpGet]
        public IActionResult GetStatus(string id)
        {
            if (Statuses.TryGetValue(id, out var status))
                return Ok(status);

            return NotFound(new { error = "Device not found", code = "DEVICE_NOT_FOUND" });
        }

        [HttpPatch]
        public IActionResult UpdateStatus(string id, [FromBody] DeviceStatus updated)
        {
            updated.DeviceId = id;
            updated.LastSeen = System.DateTime.UtcNow;
            Statuses[id] = updated;
            return Ok(updated);
        }
    }
}