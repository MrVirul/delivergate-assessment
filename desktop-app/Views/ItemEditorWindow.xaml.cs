using System;
using System.Globalization;
using System.Windows;
using DesktopApp.Models;

namespace DesktopApp.Views
{
    /// <summary>
    /// Modal editor used for both creating and updating inventory items.
    /// </summary>
    public partial class ItemEditorWindow : Window
    {
        private static readonly string[] Categories =
        {
            "flavor", "topping", "cone", "cup", "packaging", "other"
        };

        private static readonly string[] Units =
        {
            "tub", "scoop", "kg", "litre", "pcs", "box"
        };

        private readonly InventoryItem? _existingItem;

        public InventoryItem? Result { get; private set; }

        public ItemEditorWindow(InventoryItem? existingItem = null)
        {
            InitializeComponent();

            CategoryComboBox.ItemsSource = Categories;
            UnitComboBox.ItemsSource = Units;

            _existingItem = existingItem;
            if (existingItem is not null)
            {
                Title = "Edit Inventory Item";
                NameTextBox.Text = existingItem.Name;
                CategoryComboBox.SelectedItem = existingItem.Category;
                UnitComboBox.SelectedItem = existingItem.Unit;
                DescriptionTextBox.Text = existingItem.Description;
                QuantityTextBox.Text = existingItem.Quantity.ToString(CultureInfo.CurrentCulture);
                PriceTextBox.Text = (existingItem.PriceCents / 100m).ToString("0.00", CultureInfo.CurrentCulture);
                ReorderLevelTextBox.Text = existingItem.ReorderLevel.ToString(CultureInfo.CurrentCulture);
            }
            else
            {
                Title = "Add Inventory Item";
                CategoryComboBox.SelectedItem = Categories[0];
                UnitComboBox.SelectedItem = Units[0];
                QuantityTextBox.Text = "0";
                PriceTextBox.Text = "0.00";
                ReorderLevelTextBox.Text = "0";
            }

            Loaded += (_, _) => NameTextBox.Focus();
        }

        private void Save_Click(object sender, RoutedEventArgs e)
        {
            HideError();

            var name = NameTextBox.Text.Trim();
            if (name.Length is < 2 or > 120)
            {
                ShowError("Item name must be between 2 and 120 characters.");
                return;
            }

            if (CategoryComboBox.SelectedItem is not string category)
            {
                ShowError("Please choose a category.");
                return;
            }

            if (UnitComboBox.SelectedItem is not string unit)
            {
                ShowError("Please choose a unit.");
                return;
            }

            var description = DescriptionTextBox.Text.Trim();
            if (description.Length > 1000)
            {
                ShowError("Description cannot be longer than 1000 characters.");
                return;
            }

            if (!int.TryParse(QuantityTextBox.Text.Trim(), NumberStyles.Integer, CultureInfo.CurrentCulture, out var quantity) ||
                quantity < 0)
            {
                ShowError("Stock quantity must be a whole number of 0 or greater.");
                return;
            }

            if (!int.TryParse(ReorderLevelTextBox.Text.Trim(), NumberStyles.Integer, CultureInfo.CurrentCulture, out var reorderLevel) ||
                reorderLevel < 0)
            {
                ShowError("Reorder level must be a whole number of 0 or greater.");
                return;
            }

            if (!decimal.TryParse(PriceTextBox.Text.Trim(), NumberStyles.Number, CultureInfo.CurrentCulture, out var price) ||
                price < 0)
            {
                ShowError("Price must be a valid number of 0 or greater.");
                return;
            }

            Result = new InventoryItem
            {
                Id = _existingItem?.Id ?? 0,
                Name = name,
                Category = category,
                Description = description,
                Quantity = quantity,
                Unit = unit,
                PriceCents = (long)Math.Round(price * 100m, MidpointRounding.AwayFromZero),
                ReorderLevel = reorderLevel
            };

            DialogResult = true;
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
