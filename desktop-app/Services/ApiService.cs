using System;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Net.Http.Json;
using System.Text.Json;
using System.Threading.Tasks;
using DesktopApp.Models;

namespace DesktopApp.Services
{
    public class ApiService
    {
        // Ice Cream Shop API (see ../backend). All routes are under /api.
        // Configured in appsettings.json / appsettings.Local.json / environment variables.
        private static readonly string BaseUrl = AppConfig.ApiBaseUrl;
        private static readonly int PageSize = AppConfig.PageSize;

        private static readonly JsonSerializerOptions JsonOptions = new()
        {
            PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
            PropertyNameCaseInsensitive = true
        };

        private readonly HttpClient _httpClient;

        public string? Token { get; private set; }
        public string? Role { get; private set; }
        public string? CurrentUsername { get; private set; }
        public string? CurrentUserEmail { get; private set; }

        /// <summary>Human-readable reason for the most recent failed login.</summary>
        public string? LastError { get; private set; }

        public ApiService()
        {
            _httpClient = new HttpClient
            {
                BaseAddress = new Uri(BaseUrl),
                Timeout = TimeSpan.FromSeconds(30)
            };
        }

        public void SetAuthToken(string token)
        {
            Token = token;
            _httpClient.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", token);
        }

        public void ClearAuthToken()
        {
            Token = null;
            Role = null;
            CurrentUsername = null;
            CurrentUserEmail = null;
            _httpClient.DefaultRequestHeaders.Authorization = null;
        }

        public async Task<bool> LoginAsync(string email, string password)
        {
            LastError = null;

            try
            {
                var request = new LoginRequest { Email = email, Password = password };
                using var response = await _httpClient.PostAsJsonAsync("auth/login", request, JsonOptions);

                if (!response.IsSuccessStatusCode)
                {
                    var serverMessage = await ReadErrorMessageAsync(response);
                    LastError = response.StatusCode is HttpStatusCode.Unauthorized or HttpStatusCode.Forbidden
                        ? "Invalid email or password."
                        : string.IsNullOrWhiteSpace(serverMessage)
                            ? $"Login failed ({(int)response.StatusCode} {response.ReasonPhrase})."
                            : serverMessage;
                    return false;
                }

                var envelope = await ReadAsync<ApiEnvelope<LoginResponseDto>>(response);
                var data = envelope?.Data;
                if (data is null || string.IsNullOrWhiteSpace(data.AccessToken))
                {
                    LastError = "The server did not return an authentication token.";
                    return false;
                }

                SetAuthToken(data.AccessToken);
                Role = !string.IsNullOrWhiteSpace(data.User?.Role)
                    ? data.User!.Role
                    : JwtRoleReader.GetRole(data.AccessToken);
                CurrentUsername = !string.IsNullOrWhiteSpace(data.User?.Name)
                    ? data.User!.Name
                    : email;
                CurrentUserEmail = data.User?.Email ?? email;
                return true;
            }
            catch (HttpRequestException)
            {
                LastError = $"Unable to connect to the API at {BaseUrl}. Please verify the backend is running.";
                return false;
            }
            catch (TaskCanceledException)
            {
                LastError = "The login request timed out. Please try again.";
                return false;
            }
        }

        public async Task LogoutAsync()
        {
            if (string.IsNullOrWhiteSpace(Token))
            {
                ClearAuthToken();
                return;
            }

            try
            {
                using var response = await _httpClient.PostAsync("auth/logout", null);
                // The local session is cleared even if the server cannot revoke the token.
            }
            catch (HttpRequestException)
            {
                // Ignore: logging out locally is still the correct behaviour.
            }
            catch (TaskCanceledException)
            {
                // Ignore: logging out locally is still the correct behaviour.
            }
            finally
            {
                ClearAuthToken();
            }
        }

        public Task<List<InventoryItem>> GetInventoryAsync() =>
            ExecuteAsync(async () =>
            {
                var results = new List<InventoryItem>();
                var page = 1;

                while (true)
                {
                    using var response = await _httpClient.GetAsync(
                        $"items?page={page}&limit={PageSize}&sort=name");
                    await EnsureSuccessAsync(response);

                    var envelope = await ReadAsync<ApiEnvelope<List<InventoryItem>>>(response);
                    if (envelope?.Data is not { Count: > 0 })
                    {
                        break;
                    }

                    results.AddRange(envelope.Data);

                    var totalPages = envelope.Meta?.TotalPages ?? 1;
                    if (page >= totalPages)
                    {
                        break;
                    }

                    page++;
                }

                return results;
            });

        public Task<InventoryItem?> CreateItemAsync(InventoryItem item) =>
            ExecuteAsync(async () =>
            {
                using var response = await _httpClient.PostAsJsonAsync("items", ToItemRequest(item), JsonOptions);
                await EnsureSuccessAsync(response);
                return (await ReadAsync<ApiEnvelope<InventoryItem>>(response))?.Data;
            });

        public Task<InventoryItem?> UpdateItemAsync(InventoryItem item) =>
            ExecuteAsync(async () =>
            {
                using var response = await _httpClient.PutAsJsonAsync(
                    $"items/{item.Id}",
                    ToItemRequest(item),
                    JsonOptions);
                await EnsureSuccessAsync(response);
                return (await ReadAsync<ApiEnvelope<InventoryItem>>(response))?.Data;
            });

        public Task DeleteItemAsync(uint id) =>
            ExecuteAsync(async () =>
            {
                using var response = await _httpClient.DeleteAsync($"items/{id}");
                await EnsureSuccessAsync(response);
            });

        public Task<List<Customer>> GetCustomersAsync() =>
            ExecuteAsync(async () =>
            {
                var results = new List<Customer>();
                var page = 1;

                while (true)
                {
                    using var response = await _httpClient.GetAsync(
                        $"customers?page={page}&limit={PageSize}&sort=name");
                    await EnsureSuccessAsync(response);

                    var envelope = await ReadAsync<ApiEnvelope<List<Customer>>>(response);
                    if (envelope?.Data is not { Count: > 0 })
                    {
                        break;
                    }

                    results.AddRange(envelope.Data);

                    var totalPages = envelope.Meta?.TotalPages ?? 1;
                    if (page >= totalPages)
                    {
                        break;
                    }

                    page++;
                }

                return results;
            });

        public Task<Customer?> CreateCustomerAsync(Customer customer) =>
            ExecuteAsync(async () =>
            {
                using var response = await _httpClient.PostAsJsonAsync(
                    "customers",
                    ToCustomerRequest(customer),
                    JsonOptions);
                await EnsureSuccessAsync(response);
                return (await ReadAsync<ApiEnvelope<Customer>>(response))?.Data;
            });

        public Task<Customer?> UpdateCustomerAsync(Customer customer) =>
            ExecuteAsync(async () =>
            {
                using var response = await _httpClient.PutAsJsonAsync(
                    $"customers/{customer.Id}",
                    ToCustomerRequest(customer),
                    JsonOptions);
                await EnsureSuccessAsync(response);
                return (await ReadAsync<ApiEnvelope<Customer>>(response))?.Data;
            });

        public Task DeleteCustomerAsync(uint id) =>
            ExecuteAsync(async () =>
            {
                using var response = await _httpClient.DeleteAsync($"customers/{id}");
                await EnsureSuccessAsync(response);
            });

        private static ItemRequest ToItemRequest(InventoryItem item) =>
            new()
            {
                Name = item.Name,
                Category = item.Category,
                Description = item.Description,
                Quantity = item.Quantity,
                Unit = item.Unit,
                PriceCents = item.PriceCents,
                ReorderLevel = item.ReorderLevel
            };

        private static CustomerRequest ToCustomerRequest(Customer customer) =>
            new()
            {
                Name = customer.Name,
                Email = string.IsNullOrWhiteSpace(customer.Email) ? null : customer.Email.Trim(),
                Phone = customer.Phone,
                Address = customer.Address,
                Notes = customer.Notes
            };

        private static async Task<T> ExecuteAsync<T>(Func<Task<T>> action)
        {
            try
            {
                return await action();
            }
            catch (ApiException)
            {
                throw;
            }
            catch (HttpRequestException ex)
            {
                throw new ApiException(
                    $"Unable to connect to the API at {BaseUrl}. Please verify the backend is running.",
                    null,
                    ex);
            }
            catch (TaskCanceledException ex)
            {
                throw new ApiException("The request timed out. Please try again.", null, ex);
            }
        }

        private static async Task ExecuteAsync(Func<Task> action)
        {
            try
            {
                await action();
            }
            catch (ApiException)
            {
                throw;
            }
            catch (HttpRequestException ex)
            {
                throw new ApiException(
                    $"Unable to connect to the API at {BaseUrl}. Please verify the backend is running.",
                    null,
                    ex);
            }
            catch (TaskCanceledException ex)
            {
                throw new ApiException("The request timed out. Please try again.", null, ex);
            }
        }

        private static async Task EnsureSuccessAsync(HttpResponseMessage response)
        {
            if (response.IsSuccessStatusCode)
            {
                return;
            }

            var serverMessage = await ReadErrorMessageAsync(response);

            switch (response.StatusCode)
            {
                case HttpStatusCode.Forbidden:
                    throw new ApiException(
                        "You lack sufficient permissions to perform this operation.",
                        response.StatusCode);

                case HttpStatusCode.Unauthorized:
                    throw new ApiException(
                        "Your session is not authorized. Please log in again.",
                        response.StatusCode);

                default:
                    throw new ApiException(
                        string.IsNullOrWhiteSpace(serverMessage)
                            ? $"The server rejected the request ({(int)response.StatusCode})."
                            : serverMessage,
                        response.StatusCode);
            }
        }

        private static async Task<T?> ReadAsync<T>(HttpResponseMessage response)
        {
            if (response.StatusCode == HttpStatusCode.NoContent)
            {
                return default;
            }

            try
            {
                return await response.Content.ReadFromJsonAsync<T>(JsonOptions);
            }
            catch (JsonException)
            {
                return default;
            }
        }

        private static async Task<string?> ReadErrorMessageAsync(HttpResponseMessage response)
        {
            try
            {
                var body = await response.Content.ReadAsStringAsync();
                if (string.IsNullOrWhiteSpace(body))
                {
                    return null;
                }

                var envelope = JsonSerializer.Deserialize<ApiErrorEnvelope>(body, JsonOptions);
                var error = envelope?.Error;
                if (error is null)
                {
                    return body.Length > 200 ? body[..200] : body;
                }

                var message = error.Message;
                if (error.Details is { Count: > 0 })
                {
                    var details = string.Join(
                        " ",
                        error.Details.Select(detail =>
                            string.IsNullOrWhiteSpace(detail.Field)
                                ? detail.Message
                                : $"{detail.Field}: {detail.Message}"));
                    message = string.IsNullOrWhiteSpace(message) ? details : $"{message} {details}";
                }

                return message;
            }
            catch (JsonException)
            {
                return null;
            }
            catch (HttpRequestException)
            {
                return null;
            }
        }
    }
}
