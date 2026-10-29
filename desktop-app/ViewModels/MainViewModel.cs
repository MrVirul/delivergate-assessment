using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.Net;
using System.Threading.Tasks;
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using DesktopApp.Models;
using DesktopApp.Services;

namespace DesktopApp.ViewModels
{
    public class MainViewModel : ObservableObject
    {
        private const string ForbiddenMessage =
            "Action Error: You lack sufficient permissions to perform this operation.";

        private readonly ApiService _apiService;
        private readonly Func<InventoryItem?, InventoryItem?> _itemEditor;
        private readonly Func<Customer?, Customer?> _customerEditor;
        private readonly Func<string, string, bool>? _confirmDialog;

        private InventoryItem? _selectedItem;
        private Customer? _selectedCustomer;
        private string _actionError = string.Empty;
        private bool _isActionErrorVisible;
        private bool _isBusy;
        private bool _isLoggingOut;
        private string _logoutStatusMessage = string.Empty;

        public MainViewModel(
            ApiService apiService,
            Func<InventoryItem?, InventoryItem?> itemEditor,
            Func<Customer?, Customer?> customerEditor,
            Func<string, string, bool>? confirmDialog = null)
        {
            _apiService = apiService;
            _itemEditor = itemEditor;
            _customerEditor = customerEditor;
            _confirmDialog = confirmDialog;

            RefreshCommand = new AsyncRelayCommand(LoadDataAsync, () => !IsBusy && !IsLoggingOut);
            AddItemCommand = new AsyncRelayCommand(AddItemAsync, () => CanCreateItem && !IsBusy && !IsLoggingOut);
            EditItemCommand = new AsyncRelayCommand(EditItemAsync, () => CanEditItem && SelectedItem is not null && !IsBusy && !IsLoggingOut);
            DeleteItemCommand = new AsyncRelayCommand(DeleteItemAsync, () => CanDeleteItem && SelectedItem is not null && !IsBusy && !IsLoggingOut);
            AddCustomerCommand = new AsyncRelayCommand(AddCustomerAsync, () => CanModifyCustomers && !IsBusy && !IsLoggingOut);
            EditCustomerCommand = new AsyncRelayCommand(EditCustomerAsync, () => CanModifyCustomers && SelectedCustomer is not null && !IsBusy && !IsLoggingOut);
            DeleteCustomerCommand = new AsyncRelayCommand(DeleteCustomerAsync, () => CanDeleteCustomer && SelectedCustomer is not null && !IsBusy && !IsLoggingOut);
            LogoutCommand = new AsyncRelayCommand(LogoutAsync, () => CanLogout);
        }

        public ObservableCollection<InventoryItem> Inventory { get; } = new();
        public ObservableCollection<Customer> Customers { get; } = new();

        public IAsyncRelayCommand RefreshCommand { get; }
        public IAsyncRelayCommand AddItemCommand { get; }
        public IAsyncRelayCommand EditItemCommand { get; }
        public IAsyncRelayCommand DeleteItemCommand { get; }
        public IAsyncRelayCommand AddCustomerCommand { get; }
        public IAsyncRelayCommand EditCustomerCommand { get; }
        public IAsyncRelayCommand DeleteCustomerCommand { get; }
        public IAsyncRelayCommand LogoutCommand { get; }

        public string UserRole => NormalizeRole(_apiService.Role);

        public string CurrentUsername => string.IsNullOrWhiteSpace(_apiService.CurrentUsername)
            ? "Unknown"
            : _apiService.CurrentUsername;

        public string UserSummary => $"Signed in as {CurrentUsername} ({UserRole})";

        // RBAC mirrors the backend permission matrix:
        // Owner: all (create, update, delete for items and customers);
        // Manager: item create/update + customer create/update;
        // Cashier: item update + customer create/update.
        public bool CanCreateItem => (HasRole("Owner") || HasRole("Manager")) && !IsLoggingOut;

        public bool CanEditItem => (CanCreateItem || HasRole("Cashier")) && !IsLoggingOut;

        public bool CanDeleteItem => HasRole("Owner") && !IsLoggingOut;

        public bool CanModifyCustomers =>
            (HasRole("Owner") || HasRole("Manager") || HasRole("Cashier")) && !IsLoggingOut;

        public bool CanDeleteCustomer => HasRole("Owner") && !IsLoggingOut;

        public InventoryItem? SelectedItem
        {
            get => _selectedItem;
            set
            {
                if (SetProperty(ref _selectedItem, value))
                {
                    EditItemCommand.NotifyCanExecuteChanged();
                    DeleteItemCommand.NotifyCanExecuteChanged();
                }
            }
        }

        public Customer? SelectedCustomer
        {
            get => _selectedCustomer;
            set
            {
                if (SetProperty(ref _selectedCustomer, value))
                {
                    EditCustomerCommand.NotifyCanExecuteChanged();
                    DeleteCustomerCommand.NotifyCanExecuteChanged();
                }
            }
        }

        public string ActionError
        {
            get => _actionError;
            private set => SetProperty(ref _actionError, value);
        }

        public bool IsActionErrorVisible
        {
            get => _isActionErrorVisible;
            private set => SetProperty(ref _isActionErrorVisible, value);
        }

        public bool IsBusy
        {
            get => _isBusy;
            private set
            {
                if (SetProperty(ref _isBusy, value))
                {
                    NotifyAllCommandsCanExecuteChanged();
                }
            }
        }

        public bool IsLoggingOut
        {
            get => _isLoggingOut;
            private set
            {
                if (SetProperty(ref _isLoggingOut, value))
                {
                    OnPropertyChanged(nameof(CanLogout));
                    OnPropertyChanged(nameof(CanCreateItem));
                    OnPropertyChanged(nameof(CanEditItem));
                    OnPropertyChanged(nameof(CanDeleteItem));
                    OnPropertyChanged(nameof(CanModifyCustomers));
                    OnPropertyChanged(nameof(CanDeleteCustomer));
                    NotifyAllCommandsCanExecuteChanged();
                }
            }
        }

        public string LogoutStatusMessage
        {
            get => _logoutStatusMessage;
            private set => SetProperty(ref _logoutStatusMessage, value);
        }

        public bool CanLogout => !IsLoggingOut;

        /// <summary>Raised when the logout workflow has finished so the view can return to the login screen.</summary>
        public event EventHandler? LogoutCompleted;

        private void NotifyAllCommandsCanExecuteChanged()
        {
            RefreshCommand.NotifyCanExecuteChanged();
            AddItemCommand.NotifyCanExecuteChanged();
            EditItemCommand.NotifyCanExecuteChanged();
            DeleteItemCommand.NotifyCanExecuteChanged();
            AddCustomerCommand.NotifyCanExecuteChanged();
            EditCustomerCommand.NotifyCanExecuteChanged();
            DeleteCustomerCommand.NotifyCanExecuteChanged();
            LogoutCommand.NotifyCanExecuteChanged();
        }

        public async Task LogoutAsync()
        {
            if (IsLoggingOut)
            {
                return;
            }

            IsLoggingOut = true;
            LogoutStatusMessage = "Logging out, please wait...";

            try
            {
                await _apiService.LogoutAsync();
                LogoutStatusMessage = "Session ended. Returning to login...";
                await Task.Delay(400);
                LogoutCompleted?.Invoke(this, EventArgs.Empty);
            }
            finally
            {
                IsLoggingOut = false;
            }
        }

        public async Task LoadDataAsync()
        {
            if (IsBusy)
            {
                return;
            }

            ClearError();
            IsBusy = true;

            try
            {
                var inventoryTask = _apiService.GetInventoryAsync();
                var customersTask = _apiService.GetCustomersAsync();

                await Task.WhenAll(inventoryTask, customersTask);

                ReplaceAll(Inventory, await inventoryTask);
                ReplaceAll(Customers, await customersTask);
            }
            catch (ApiException ex)
            {
                ShowError(ex);
            }
            catch (Exception ex)
            {
                ShowError($"Unexpected error: {ex.Message}");
            }
            finally
            {
                IsBusy = false;
            }
        }

        private async Task AddItemAsync()
        {
            var draft = _itemEditor(null);
            if (draft is null)
            {
                return;
            }

            await ExecuteActionAsync(async () =>
            {
                var created = await _apiService.CreateItemAsync(draft);
                Inventory.Add(created ?? draft);
            });
        }

        private async Task EditItemAsync()
        {
            var original = SelectedItem;
            if (original is null)
            {
                return;
            }

            var edited = _itemEditor(original);
            if (edited is null)
            {
                return;
            }

            await ExecuteActionAsync(async () =>
            {
                var updated = await _apiService.UpdateItemAsync(edited) ?? edited;
                var index = Inventory.IndexOf(original);
                if (index >= 0)
                {
                    Inventory[index] = updated;
                }
            });
        }

        private async Task DeleteItemAsync()
        {
            var item = SelectedItem;
            if (item is null)
            {
                return;
            }

            if (!CanDeleteItem)
            {
                ShowError(ForbiddenMessage);
                return;
            }

            var confirmed = _confirmDialog?.Invoke(
                $"Are you sure you want to delete '{item.Name}' (ID: {item.Id})?",
                "Confirm Delete Item") ?? true;

            if (!confirmed)
            {
                return;
            }

            await ExecuteActionAsync(async () =>
            {
                await _apiService.DeleteItemAsync(item.Id);
                Inventory.Remove(item);
                SelectedItem = null;
            });
        }

        private async Task AddCustomerAsync()
        {
            var draft = _customerEditor(null);
            if (draft is null)
            {
                return;
            }

            await ExecuteActionAsync(async () =>
            {
                var created = await _apiService.CreateCustomerAsync(draft);
                Customers.Add(created ?? draft);
            });
        }

        private async Task EditCustomerAsync()
        {
            var original = SelectedCustomer;
            if (original is null)
            {
                return;
            }

            var edited = _customerEditor(original);
            if (edited is null)
            {
                return;
            }

            await ExecuteActionAsync(async () =>
            {
                var updated = await _apiService.UpdateCustomerAsync(edited) ?? edited;
                var index = Customers.IndexOf(original);
                if (index >= 0)
                {
                    Customers[index] = updated;
                }
            });
        }

        private async Task DeleteCustomerAsync()
        {
            var customer = SelectedCustomer;
            if (customer is null)
            {
                return;
            }

            if (!CanDeleteCustomer)
            {
                ShowError(ForbiddenMessage);
                return;
            }

            var confirmed = _confirmDialog?.Invoke(
                $"Are you sure you want to delete customer '{customer.Name}' (ID: {customer.Id})?",
                "Confirm Delete Customer") ?? true;

            if (!confirmed)
            {
                return;
            }

            await ExecuteActionAsync(async () =>
            {
                await _apiService.DeleteCustomerAsync(customer.Id);
                Customers.Remove(customer);
                SelectedCustomer = null;
            });
        }

        private async Task ExecuteActionAsync(Func<Task> action)
        {
            ClearError();
            try
            {
                await action();
            }
            catch (ApiException ex)
            {
                ShowError(ex);
            }
            catch (Exception ex)
            {
                ShowError($"Unexpected error: {ex.Message}");
            }
        }

        private void ShowError(ApiException exception)
        {
            ShowError(exception.StatusCode == HttpStatusCode.Forbidden
                ? ForbiddenMessage
                : exception.Message);
        }

        private void ShowError(string message)
        {
            ActionError = message;
            IsActionErrorVisible = true;
        }

        private void ClearError()
        {
            ActionError = string.Empty;
            IsActionErrorVisible = false;
        }

        private static void ReplaceAll<T>(ObservableCollection<T> target, IEnumerable<T> source)
        {
            target.Clear();
            foreach (var item in source)
            {
                target.Add(item);
            }
        }

        private bool HasRole(string role) =>
            string.Equals(UserRole, role, StringComparison.OrdinalIgnoreCase);

        private static string NormalizeRole(string? role)
        {
            if (string.IsNullOrWhiteSpace(role))
            {
                return "Unknown";
            }

            var normalized = role.Trim();
            if (normalized.StartsWith("ROLE_", StringComparison.OrdinalIgnoreCase))
            {
                normalized = normalized[5..];
            }

            return normalized;
        }
    }
}
