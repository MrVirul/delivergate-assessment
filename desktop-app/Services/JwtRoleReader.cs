using System;
using System.Text;
using System.Text.Json;

namespace DesktopApp.Services
{
    /// <summary>
    /// Reads the role claim from a JWT when the login response does not include one.
    /// </summary>
    public static class JwtRoleReader
    {
        private static readonly string[] RoleClaimNames =
        {
            "role",
            "roles",
            "http://schemas.microsoft.com/ws/2008/06/identity/claims/role"
        };

        public static string? GetRole(string token)
        {
            try
            {
                var parts = token.Split('.');
                if (parts.Length < 2)
                {
                    return null;
                }

                var payload = parts[1].Replace('-', '+').Replace('_', '/');
                payload = payload.PadRight(payload.Length + (4 - payload.Length % 4) % 4, '=');

                var json = Encoding.UTF8.GetString(Convert.FromBase64String(payload));
                using var document = JsonDocument.Parse(json);

                foreach (var name in RoleClaimNames)
                {
                    if (!document.RootElement.TryGetProperty(name, out var claim))
                    {
                        continue;
                    }

                    if (claim.ValueKind == JsonValueKind.String)
                    {
                        return claim.GetString();
                    }

                    if (claim.ValueKind == JsonValueKind.Array &&
                        claim.GetArrayLength() > 0 &&
                        claim[0].ValueKind == JsonValueKind.String)
                    {
                        return claim[0].GetString();
                    }
                }
            }
            catch (Exception)
            {
                // A malformed token simply means we cannot determine the role.
            }

            return null;
        }
    }
}
