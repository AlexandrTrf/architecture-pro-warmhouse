using SmartHomeApi.Controllers.Requests;
using SmartHomeApi.Models;

namespace SmartHomeApi.Mapper;

public static class Mapper
{
    public static Device ToDevice(this DeviceRequest request)
    {
        return new Device
        {
            Name = request.Name,
            Type = request.Type,
            Location = request.Location,
            Manufacturer = request.Manufacturer,
            Model = request.Model,
            IpAddress = request.IpAddress,
            MacAddress = request.MacAddress,
            Status = "online",
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };
    }
}