using System;
using System.Threading.Tasks;
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using DesktopApp.Services;

namespace DesktopApp.ViewModels
{
    public partial class LoginViewModel : ObservableObject
    {
        private readonly ApiService _apiService;

        [ObservableProperty]
        private string email = string.Empty;

        [ObservableProperty]
        private string password = string.Empty;

        [ObservableProperty]
        private string errorMessage = string.Empty;

        [ObservableProperty]
        private bool isErrorVisible;

        [ObservableProperty]
        private string infoMessage = string.Empty;

        [ObservableProperty]
        private bool isInfoVisible;

        [ObservableProperty]
        private bool isBusy;

        public IAsyncRelayCommand LoginCommand { get; }

        /// <summary>Raised after a successful authentication so the view can open the dashboard.</summary>
        public event EventHandler? LoginSucceeded;

        public LoginViewModel(ApiService apiService)
        {
            _apiService = apiService;
            LoginCommand = new AsyncRelayCommand(ExecuteLoginAsync);
        }

        public void ShowInfo(string message)
        {
            InfoMessage = message;
            IsInfoVisible = true;
            ErrorMessage = string.Empty;
            IsErrorVisible = false;
        }

        private async Task ExecuteLoginAsync()
        {
            IsErrorVisible = false;
            ErrorMessage = string.Empty;
            IsInfoVisible = false;
            InfoMessage = string.Empty;

            if (string.IsNullOrWhiteSpace(Email) || string.IsNullOrWhiteSpace(Password))
            {
                ShowError("Please enter both email and password.");
                return;
            }

            IsBusy = true;
            try
            {
                var success = await _apiService.LoginAsync(Email.Trim(), Password);
                if (success)
                {
                    Password = string.Empty;
                    LoginSucceeded?.Invoke(this, EventArgs.Empty);
                }
                else
                {
                    ShowError(string.IsNullOrWhiteSpace(_apiService.LastError)
                        ? "Invalid email or password."
                        : _apiService.LastError);
                }
            }
            catch (Exception ex)
            {
                ShowError($"Login failed: {ex.Message}");
            }
            finally
            {
                IsBusy = false;
            }
        }

        private void ShowError(string message)
        {
            ErrorMessage = message;
            IsErrorVisible = true;
        }
    }
}
