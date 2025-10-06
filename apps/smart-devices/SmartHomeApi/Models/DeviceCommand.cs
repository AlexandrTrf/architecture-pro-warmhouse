using System.Collections.Generic;

namespace SmartHomeApi.Models
{
    public class DeviceCommand
    {
        public string Command { get; set; } = default!;
        public Dictionary<string, object>? Parameters { get; set; }
        public string? Priority { get; set; }
    }
}