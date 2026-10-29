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
            _viewModel = new MainViewModel(apiService, ShowItemEditor, ShowCustomerEditor);
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

        private async void Logout_Click(object sender, RoutedEventArgs e)
        {
            await _apiService.LogoutAsync();

            var loginWindow = new LoginWindow();
            Application.Current.MainWindow = loginWindow;
            loginWindow.Show();
            Close();
        }
    }
}
