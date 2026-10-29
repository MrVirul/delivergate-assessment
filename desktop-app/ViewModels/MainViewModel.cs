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

        private InventoryItem? _selectedItem;
        private Customer? _selectedCustomer;
        private string _actionError = string.Empty;
        private bool _isActionErrorVisible;
        private bool _isBusy;

        public MainViewModel(
            ApiService apiService,
            Func<InventoryItem?, InventoryItem?> itemEditor,
            Func<Customer?, Customer?> customerEditor)
        {
            _apiService = apiService;
            _itemEditor = itemEditor;
            _customerEditor = customerEditor;

            RefreshCommand = new AsyncRelayCommand(LoadDataAsync, () => !IsBusy);
            AddItemCommand = new AsyncRelayCommand(AddItemAsync, () => CanCreateItem);
            EditItemCommand = new AsyncRelayCommand(EditItemAsync, () => CanEditItem && SelectedItem is not null);
            AddCustomerCommand = new AsyncRelayCommand(AddCustomerAsync, () => CanModifyCustomers);
            EditCustomerCommand = new AsyncRelayCommand(EditCustomerAsync, () => CanModifyCustomers && SelectedCustomer is not null);
        }

        public ObservableCollection<InventoryItem> Inventory { get; } = new();
        public ObservableCollection<Customer> Customers { get; } = new();

        public IAsyncRelayCommand RefreshCommand { get; }
        public IAsyncRelayCommand AddItemCommand { get; }
        public IAsyncRelayCommand EditItemCommand { get; }
        public IAsyncRelayCommand AddCustomerCommand { get; }
        public IAsyncRelayCommand EditCustomerCommand { get; }

        public string UserRole => NormalizeRole(_apiService.Role);

        public string CurrentUsername => string.IsNullOrWhiteSpace(_apiService.CurrentUsername)
            ? "Unknown"
            : _apiService.CurrentUsername;

        public string UserSummary => $"Signed in as {CurrentUsername} ({UserRole})";

        // RBAC mirrors the backend permission matrix:
        // Owner: all; Manager: item create/update + customer create/update;
        // Cashier: item update + customer create/update.
        public bool CanCreateItem => HasRole("Owner") || HasRole("Manager");

        public bool CanEditItem => CanCreateItem || HasRole("Cashier");

        public bool CanModifyCustomers =>
            HasRole("Owner") || HasRole("Manager") || HasRole("Cashier");

        public InventoryItem? SelectedItem
        {
            get => _selectedItem;
            set
            {
                if (SetProperty(ref _selectedItem, value))
                {
                    EditItemCommand.NotifyCanExecuteChanged();
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
                    RefreshCommand.NotifyCanExecuteChanged();
                }
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
