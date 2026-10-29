using System.Collections.Generic;

namespace DesktopApp.Models
{
    /// <summary>Envelope used by every successful API response.</summary>
    public class ApiEnvelope<T>
    {
        public T? Data { get; set; }
        public ListMeta? Meta { get; set; }
    }

    public class ListMeta
    {
        public int Page { get; set; }
        public int Limit { get; set; }
        public long Total { get; set; }
        public int TotalPages { get; set; }
    }

    public class ApiErrorEnvelope
    {
        public ApiError? Error { get; set; }
    }

    public class ApiError
    {
        public string? Code { get; set; }
        public string? Message { get; set; }
        public List<ApiFieldError>? Details { get; set; }
    }

    public class ApiFieldError
    {
        public string? Field { get; set; }
        public string? Message { get; set; }
    }

    public class LoginRequest
    {
        public string Email { get; set; } = string.Empty;
        public string Password { get; set; } = string.Empty;
    }

    public class LoginResponseDto
    {
        public string? AccessToken { get; set; }
        public string? TokenType { get; set; }
        public int ExpiresIn { get; set; }
        public UserDto? User { get; set; }
    }

    public class UserDto
    {
        public uint Id { get; set; }
        public string? Name { get; set; }
        public string? Email { get; set; }
        public string? Role { get; set; }
    }

    public class ItemRequest
    {
        public string Name { get; set; } = string.Empty;
        public string Category { get; set; } = string.Empty;
        public string Description { get; set; } = string.Empty;
        public int Quantity { get; set; }
        public string Unit { get; set; } = string.Empty;
        public long PriceCents { get; set; }
        public int ReorderLevel { get; set; }
    }

    public class CustomerRequest
    {
        public string Name { get; set; } = string.Empty;
        public string? Email { get; set; }
        public string Phone { get; set; } = string.Empty;
        public string Address { get; set; } = string.Empty;
        public string Notes { get; set; } = string.Empty;
    }
}
