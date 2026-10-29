using System.Text.Json.Serialization;

namespace DesktopApp.Models
{
    public class InventoryItem
    {
        public uint Id { get; set; }
        public string Name { get; set; } = string.Empty;
        public string Category { get; set; } = string.Empty;
        public string Description { get; set; } = string.Empty;
        public int Quantity { get; set; }
        public string Unit { get; set; } = string.Empty;
        public long PriceCents { get; set; }
        public int ReorderLevel { get; set; }
        public bool LowStock { get; set; }

        [JsonIgnore]
        public decimal Price => PriceCents / 100m;

        [JsonIgnore]
        public string StockStatus => LowStock ? "Low stock" : "OK";
    }
}
