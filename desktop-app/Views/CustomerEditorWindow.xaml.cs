using System;
using System.Windows;
using DesktopApp.Models;

namespace DesktopApp.Views
{
    /// <summary>
    /// Modal editor used for both creating and updating customers.
    /// </summary>
    public partial class CustomerEditorWindow : Window
    {
        private readonly Customer? _existingCustomer;

        public Customer? Result { get; private set; }

        public CustomerEditorWindow(Customer? existingCustomer = null)
        {
            InitializeComponent();

            _existingCustomer = existingCustomer;
            if (existingCustomer is not null)
            {
                Title = "Edit Customer";
                FullNameTextBox.Text = existingCustomer.Name;
                EmailTextBox.Text = existingCustomer.Email ?? string.Empty;
                PhoneTextBox.Text = existingCustomer.Phone;
                AddressTextBox.Text = existingCustomer.Address;
                NotesTextBox.Text = existingCustomer.Notes;
            }
            else
            {
                Title = "Add Customer";
            }

            Loaded += (_, _) => FullNameTextBox.Focus();
        }

        private void Save_Click(object sender, RoutedEventArgs e)
        {
            HideError();

            var fullName = FullNameTextBox.Text.Trim();
            if (fullName.Length is < 2 or > 120)
            {
                ShowError("Full name must be between 2 and 120 characters.");
                return;
            }

            var email = EmailTextBox.Text.Trim();
            if (!string.IsNullOrWhiteSpace(email) &&
                (!email.Contains('@') || !email.Contains('.')))
            {
                ShowError("Please enter a valid email address.");
                return;
            }

            var phone = PhoneTextBox.Text.Trim();
            if (!IsValidPhone(phone))
            {
                ShowError("Phone must be 7-30 characters and may contain only digits, +, -, ( ), or spaces.");
                return;
            }

            var address = AddressTextBox.Text.Trim();
            if (address.Length > 255)
            {
                ShowError("Address cannot be longer than 255 characters.");
                return;
            }

            var notes = NotesTextBox.Text.Trim();
            if (notes.Length > 1000)
            {
                ShowError("Notes cannot be longer than 1000 characters.");
                return;
            }

            Result = new Customer
            {
                Id = _existingCustomer?.Id ?? 0,
                Name = fullName,
                Email = string.IsNullOrWhiteSpace(email) ? null : email,
                Phone = phone,
                Address = address,
                Notes = notes
            };

            DialogResult = true;
        }

        private static bool IsValidPhone(string phone)
        {
            if (phone.Length is < 7 or > 30)
            {
                return false;
            }

            foreach (var character in phone)
            {
                if ((character >= '0' && character <= '9') ||
                    character is '+' or '-' or '(' or ')' or ' ')
                {
                    continue;
                }

                return false;
            }

            return true;
        }

        private void ShowError(string message)
        {
            ErrorText.Text = message;
            ErrorText.Visibility = Visibility.Visible;
        }

        private void HideError()
        {
            ErrorText.Text = string.Empty;
            ErrorText.Visibility = Visibility.Collapsed;
        }
    }
}
