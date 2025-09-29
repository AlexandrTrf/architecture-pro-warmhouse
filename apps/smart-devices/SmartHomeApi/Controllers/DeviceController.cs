using Microsoft.AspNetCore.Mvc;
using SmartHomeApi.Models;
using System.Collections.Generic;
using System.Linq;
using SmartHomeApi.Controllers.Requests;
using SmartHomeApi.Mapper;

namespace SmartHomeApi.Controllers
{
    [ApiController]
    [Route("api/v1/devices")]
    public class DevicesController : ControllerBase
    {
        private static readonly List<Device> Devices = new();

        [HttpGet]
        public IActionResult GetDevices([FromQuery] string? type, [FromQuery] string? status, [FromQuery] string? location)
        {
            var query = Devices.AsQueryable();
            if (!string.IsNullOrEmpty(type))
                query = query.Where(d => d.Type == type);
            if (!string.IsNullOrEmpty(status))
                query = query.Where(d => d.Status == status);
            if (!string.IsNullOrEmpty(location))
                query = query.Where(d => d.Location == location);

            return Ok(query.ToList());
        }

        [HttpPost]
        public IActionResult CreateDevice([FromBody] DeviceRequest device)
        {
            var newDevice = device.ToDevice();
            newDevice.Id = $"dev_{Devices.Count + 1}";
            Devices.Add(newDevice);
            return CreatedAtAction(nameof(GetDeviceById), new { id = newDevice.Id }, device);
        }

        [HttpGet("{id}")]
        public IActionResult GetDeviceById(string id)
        {
            var device = Devices.FirstOrDefault(d => d.Id == id);
            if (device == null)
                return NotFound(new { error = "Device not found", code = "DEVICE_NOT_FOUND" });

            return Ok(device);
        }

        [HttpPut("{id}")]
        public IActionResult UpdateDevice(string id, [FromBody] DeviceRequest updated)
        {
            var device = Devices.FirstOrDefault(d => d.Id == id);
            if (device == null)
                return NotFound(new { error = "Device not found", code = "DEVICE_NOT_FOUND" });

            device.Name = updated.Name ?? device.Name;
            device.Type = updated.Type ?? device.Type;
            device.Location = updated.Location ?? device.Location;
            device.Status = updated.Status ?? device.Status;
            device.UpdatedAt = System.DateTime.UtcNow;

            return Ok(device);
        }

        [HttpDelete("{id}")]
        public IActionResult DeleteDevice(string id)
        {
            var device = Devices.FirstOrDefault(d => d.Id == id);
            if (device == null)
                return NotFound(new { error = "Device not found", code = "DEVICE_NOT_FOUND" });

            Devices.Remove(device);
            return Ok(new { message = "Device deleted successfully" });
        }
    }
}
