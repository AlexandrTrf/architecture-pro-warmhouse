using System;

namespace SmartHomeApi.Models
{
    public class Device
    {
        public string Id { get; set; }= default!;
        public string Name { get; set; } = default!;
        public string Type { get; set; } = default!;
        public string Location { get; set; } = default!;
        public string Status { get; set; }
        public string? Manufacturer { get; set; }
        public string? Model { get; set; }
        public string? FirmwareVersion { get; set; }
        public string? IpAddress { get; set; }
        public string? MacAddress { get; set; }
        public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
        public DateTime UpdatedAt { get; set; } = DateTime.UtcNow;
    }
}