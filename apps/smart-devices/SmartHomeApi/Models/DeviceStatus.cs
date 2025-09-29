using System;

namespace SmartHomeApi.Models
{
    public class DeviceStatus
    {
        public string DeviceId { get; set; } = default!;
        public string Status { get; set; } = default!;
        public DateTime LastSeen { get; set; } = DateTime.UtcNow;
        public string? Message { get; set; }
    }
}