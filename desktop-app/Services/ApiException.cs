using System;
using System.Net;

namespace DesktopApp.Services
{
    /// <summary>
    /// Raised when the API returns a non-success status code or cannot be reached.
    /// </summary>
    public class ApiException : Exception
    {
        public HttpStatusCode? StatusCode { get; }

        public ApiException(string message, HttpStatusCode? statusCode = null, Exception? innerException = null)
            : base(message, innerException)
        {
            StatusCode = statusCode;
        }
    }
}
