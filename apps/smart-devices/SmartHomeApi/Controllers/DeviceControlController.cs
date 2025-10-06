using Microsoft.AspNetCore.Mvc;
using SmartHomeApi.Models;
using System.Collections.Generic;

namespace SmartHomeApi.Controllers
{
    [ApiController]
    [Route("api/v1/devices/{id}/control")]
    public class DeviceControlController : ControllerBase
    {
        private static readonly Dictionary<string, List<DeviceCommand>> CommandHistory = new();

        [HttpPost]
        public IActionResult SendCommand(string id, [FromBody] DeviceCommand command)
        {
            if (!CommandHistory.ContainsKey(id))
                CommandHistory[id] = new List<DeviceCommand>();

            CommandHistory[id].Add(command);

            return Ok(new { message = "Command sent successfully", command_id = $"cmd_{CommandHistory[id].Count}" });
        }
    }
}