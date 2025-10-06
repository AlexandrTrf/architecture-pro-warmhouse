using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc.ModelBinding;
using Microsoft.Extensions.Logging;
using System.Text.Json;
using System.Threading.Tasks;

public class ValidationMiddleware
{
    private readonly RequestDelegate _next;
    private readonly ILogger<ValidationMiddleware> _logger;

    public ValidationMiddleware(RequestDelegate next, ILogger<ValidationMiddleware> logger)
    {
        _next = next;
        _logger = logger;
    }

    public async Task InvokeAsync(HttpContext context)
    {
        // Буферизация тела запроса, чтобы можно было читать его несколько раз (опционально)
        context.Request.EnableBuffering();

        // Продолжаем выполнение конвейера
        await _next(context);

        // После выполнения контроллера проверяем ModelState
        if (context.Items.ContainsKey("ModelState") && context.Items["ModelState"] is ModelStateDictionary modelState)
        {
            if (!modelState.IsValid)
            {
                var errors = modelState.Values
                    .SelectMany(v => v.Errors)
                    .Select(e => e.ErrorMessage)
                    .ToList();

                _logger.LogWarning("Model validation failed: {Errors}", string.Join("; ", errors));

                context.Response.StatusCode = StatusCodes.Status400BadRequest;
                context.Response.ContentType = "application/json";

                var response = JsonSerializer.Serialize(new
                {
                    Message = "Validation failed",
                    Errors = errors
                });

                await context.Response.WriteAsync(response);
            }
        }
    }
}