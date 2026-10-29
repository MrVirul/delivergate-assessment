using System;
using System.Windows;
using DesktopApp.Services;
using DesktopApp.ViewModels;

namespace DesktopApp.Views
{
    /// <summary>
    /// Interaction logic for LoginWindow.xaml
    /// </summary>
    public partial class LoginWindow : Window
    {
        private readonly ApiService _apiService = new();
        private readonly LoginViewModel _viewModel;

        // Parameterless constructor required by WPF when the window is created from XAML (StartupUri).
        public LoginWindow()
            : this(null)
        {
        }

        public LoginWindow(string? infoMessage)
        {
            InitializeComponent();

            _viewModel = new LoginViewModel(_apiService);
            if (!string.IsNullOrWhiteSpace(infoMessage))
            {
                _viewModel.ShowInfo(infoMessage);
            }
            _viewModel.LoginSucceeded += OnLoginSucceeded;
            DataContext = _viewModel;

            Loaded += (_, _) => EmailTextBox.Focus();
        }

        private void PasswordBox_PasswordChanged(object sender, RoutedEventArgs e)
        {
            if (DataContext is LoginViewModel viewModel)
            {
                viewModel.Password = PasswordBox.Password;
            }
        }

        private void OnLoginSucceeded(object? sender, EventArgs e)
        {
            var mainWindow = new MainWindow(_apiService);
            Application.Current.MainWindow = mainWindow;
            mainWindow.Show();
            Close();
        }
    }
}
