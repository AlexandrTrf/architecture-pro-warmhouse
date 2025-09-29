using System;
using System.Collections.Generic;

namespace SmartHomeApi.Models
{
    public class DeviceConfig
    {
        public string DeviceId { get; set; } = default!;
        public Dictionary<string, object> Settings { get; set; } = new();
        public string? FirmwareVersion { get; set; }
        public bool UpdateAvailable { get; set; }
        public DateTime LastUpdated { get; set; } = DateTime.UtcNow;
    }
}