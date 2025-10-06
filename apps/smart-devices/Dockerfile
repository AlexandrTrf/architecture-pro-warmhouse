# Build stage
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

COPY SmartHomeApi/*.csproj ./SmartHomeApi/
RUN dotnet restore SmartHomeApi/SmartHomeApi.csproj

COPY SmartHomeApi/. ./SmartHomeApi/
WORKDIR /src/SmartHomeApi
RUN dotnet publish -c Release -o /app/publish

# Runtime stage
FROM mcr.microsoft.com/dotnet/aspnet:8.0
WORKDIR /app
COPY --from=build /app/publish .
ENTRYPOINT ["dotnet", "SmartHomeApi.dll"]
