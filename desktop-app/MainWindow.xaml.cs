using System.Windows;
using DesktopApp.Models;
using DesktopApp.Services;
using DesktopApp.ViewModels;
using DesktopApp.Views;

namespace DesktopApp
{
    /// <summary>
    /// Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class MainWindow : Window
    {
        private readonly ApiService _apiService;
        private readonly MainViewModel _viewModel;

        public MainWindow(ApiService apiService)
        {
            InitializeComponent();

            _apiService = apiService;
            _viewModel = new MainViewModel(
                apiService,
                ShowItemEditor,
                ShowCustomerEditor,
                ConfirmDelete);
            _viewModel.LogoutCompleted += OnLogoutCompleted;
            DataContext = _viewModel;

            Loaded += async (_, _) => await _viewModel.LoadDataAsync();
        }

        private InventoryItem? ShowItemEditor(InventoryItem? item)
        {
            var editor = new ItemEditorWindow(item) { Owner = this };
            return editor.ShowDialog() == true ? editor.Result : null;
        }

        private Customer? ShowCustomerEditor(Customer? customer)
        {
            var editor = new CustomerEditorWindow(customer) { Owner = this };
            return editor.ShowDialog() == true ? editor.Result : null;
        }

        private bool ConfirmDelete(string message, string caption)
        {
            return MessageBox.Show(
                this,
                message,
                caption,
                MessageBoxButton.YesNo,
                MessageBoxImage.Warning) == MessageBoxResult.Yes;
        }

        private void OnLogoutCompleted(object? sender, System.EventArgs e)
        {
            var loginWindow = new LoginWindow("You have been logged out successfully.");
            Application.Current.MainWindow = loginWindow;
            loginWindow.Show();
            Close();
        }
    }
}
