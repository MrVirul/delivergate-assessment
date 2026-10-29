using System;
using Microsoft.Extensions.Configuration;

namespace DesktopApp.Services
{
    /// <summary>
    /// Application settings loaded from appsettings.json, then the optional
    /// (gitignored) appsettings.Local.json, then environment variables.
    /// </summary>
    internal static class AppConfig
    {
        private const string DefaultApiBaseUrl = "http://localhost:8080/api/";
        private const int DefaultPageSize = 100;

        private static readonly IConfigurationRoot Configuration = new ConfigurationBuilder()
            .SetBasePath(AppContext.BaseDirectory)
            .AddJsonFile("appsettings.json", optional: true, reloadOnChange: false)
            .AddJsonFile("appsettings.Local.json", optional: true, reloadOnChange: false)
            .AddEnvironmentVariables(prefix: "DESKTOPAPP_")
            .Build();

        /// <summary>
        /// API root, always with a trailing slash. Override with Api:BaseUrl in
        /// appsettings.Local.json or the DESKTOPAPP_Api__BaseUrl environment variable.
        /// </summary>
        public static string ApiBaseUrl { get; } = NormalizeBaseUrl(Configuration["Api:BaseUrl"]);

        /// <summary>
        /// Number of records requested per page. Override with Api:PageSize or
        /// the DESKTOPAPP_Api__PageSize environment variable.
        /// </summary>
        public static int PageSize { get; } = ReadPageSize();

        private static int ReadPageSize() =>
            int.TryParse(Configuration["Api:PageSize"], out var pageSize) && pageSize > 0
                ? pageSize
                : DefaultPageSize;

        private static string NormalizeBaseUrl(string? value)
        {
            var url = string.IsNullOrWhiteSpace(value) ? DefaultApiBaseUrl : value.Trim();
            return url.EndsWith('/') ? url : url + "/";
        }
    }
}
