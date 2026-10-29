Screen 1 - The Login Screen
The login window handles user credentials, calls the ApiService, manages validation/error states, and transitions to the main dashboard upon success.

Step 1: Create LoginViewModel.cs
In your ViewModels folder, create LoginViewModel.cs inheriting from ObservableObject:

C#
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using DesktopApp.Services;
using System.Threading.Tasks;

public partial class LoginViewModel : ObservableObject
{
    private readonly ApiService _apiService;

    [ObservableProperty]
    private string username;

    [ObservableProperty]
    private string password;

    [ObservableProperty]
    private string errorMessage;

    [ObservableProperty]
    private bool isErrorVisible;

    public IAsyncRelayCommand LoginCommand { get; }

    public LoginViewModel(ApiService apiService)
    {
        _apiService = apiService;
        LoginCommand = new AsyncRelayCommand(ExecuteLoginAsync);
    }

    private async Task ExecuteLoginAsync()
    {
        IsErrorVisible = false;
        
        if (string.IsNullOrWhiteSpace(Username) || string.IsNullOrWhiteSpace(Password))
        {
            ErrorMessage = "Please enter both username and password.";
            IsErrorVisible = true;
            return;
        }

        bool success = await _apiService.LoginAsync(Username, Password);
        if (success)
        {
            // Open Main Dashboard and close Login Window
            var mainWindow = new MainWindow(_apiService);
            mainWindow.Show();
            App.Current.Windows[0]?.Close();
        }
        else
        {
            ErrorMessage = "Invalid credentials or unauthorized access.";
            IsErrorVisible = true;
        }
    }
}
Step 2: Design LoginWindow.xaml
Update your primary startup window or create a dedicated LoginWindow.xaml:

XML
<Window x:Class="DesktopApp.Views.LoginWindow"
        xmlns="http://schemas.microsoft.com/winfx/2006/xaml/presentation"
        xmlns:x="http://schemas.microsoft.com/winfx/2006/xaml"
        Title="System Login" Height="350" Width="400" WindowStartupLocation="CenterScreen">
    <Grid Margin="20">
        <StackPanel VerticalAlignment="Center">
            <TextBlock Text="Welcome Back" FontSize="22" FontWeight="Bold" Margin="0,0,0,20"/>
            
            <TextBlock Text="Username" FontWeight="SemiBold"/>
            <TextBox Text="{Binding Username, UpdateSourceTrigger=PropertyChanged}" Height="30" Margin="0,5,0,15"/>
            
            <TextBlock Text="Password" FontWeight="SemiBold"/>
            <PasswordBox x:Name="PasswordBox" Height="30" Margin="0,5,0,15"/>

            <Button Content="Login" Command="{Binding LoginCommand}" Height="35" Background="#007ACC" Foreground="White" FontWeight="Bold"/>

            <TextBlock Text="{Binding ErrorMessage}" Foreground="Red" Margin="0,15,0,0" 
                       Visibility="{Binding IsErrorVisible, Converter={StaticResource BooleanToVisibilityConverter}}"/>
        </StackPanel>
    </Grid>
</Window>
Phase 4: Screen 2 & 3 - Dashboard with Inventory & Customer Management
The main dashboard (MainWindow.xaml) typically utilizes a sidebar or a TabControl to toggle between Inventory Management and Customer Management.

Step 1: Design MainWindow.xaml
Use a layout with navigation on the left or top tabs to switch contexts easily:

XML
<Window x:Class="DesktopApp.MainWindow"
        xmlns="http://schemas.microsoft.com/winfx/2006/xaml/presentation"
        xmlns:x="http://schemas.microsoft.com/winfx/2006/xaml"
        Title="Enterprise Management System" Height="600" Width="900" WindowStartupLocation="CenterScreen">
    <TabControl>
        <!-- Inventory Tab -->
        <TabItem Header="Inventory Management">
            <Grid Margin="10">
                <Grid.RowDefinitions>
                    <RowDefinition Height="Auto"/>
                    <RowDefinition Height="*"/>
                    <RowDefinition Height="Auto"/>
                </Grid.RowDefinitions>
                
                <StackPanel Orientation="Horizontal" Margin="0,0,0,10">
                    <Button Content="Add Item" Width="90" Height="30" Margin="0,0,10,0"/>
                    <Button Content="Edit Selected" Width="90" Height="30"/>
                </StackPanel>

                <DataGrid Grid.Row="1" AutoGenerateColumns="False" CanUserAddRows="False">
                    <DataGrid.Columns>
                        <DataGridTextColumn Header="Item ID" Width="100"/>
                        <DataGridTextColumn Header="Item Name" Width="250"/>
                        <DataGridTextColumn Header="Stock Quantity" Width="120"/>
                        <DataGridTextColumn Header="Price ($)" Width="100"/>
                    </DataGrid.Columns>
                </DataGrid>
                
                <TextBlock Grid.Row="2" Text="* Restricted permissions: Only Managers can edit inventory items." Foreground="Gray" Margin="0,10,0,0"/>
            </Grid>
        </TabItem>

        <!-- Customer Tab -->
        <TabItem Header="Customer Management">
            <Grid Margin="10">
                <Grid.RowDefinitions>
                    <RowDefinition Height="Auto"/>
                    <RowDefinition Height="*"/>
                </Grid.RowDefinitions>

                <StackPanel Orientation="Horizontal" Margin="0,0,0,10">
                    <Button Content="Add Customer" Width="110" Height="30" Margin="0,0,10,0"/>
                    <Button Content="Edit Customer" Width="110" Height="30"/>
                </StackPanel>

                <DataGrid Grid.Row="1" AutoGenerateColumns="False" CanUserAddRows="False">
                    <DataGrid.Columns>
                        <DataGridTextColumn Header="Customer ID" Width="100"/>
                        <DataGridTextColumn Header="Full Name" Width="250"/>
                        <DataGridTextColumn Header="Email" Width="200"/>
                        <DataGridTextColumn Header="Phone Number" Width="150"/>
                    </DataGrid.Columns>
                </DataGrid>
            </Grid>
        </TabItem>
    </TabControl>
</Window>
Phase 5: Handling Permissions & Action Errors
To handle role-based restrictions gracefully (e.g., Cashiers can view items, but only Admins can add/edit):

Decode the JWT token on login or store the Role returned in your LoginResponseDto.

Bind the IsEnabled property of your action buttons (Add/Edit) to a boolean property in your ViewModel (e.g., CanModifyData).

C#
// Example ViewModel logic snippet for permission checks
public bool CanModifyData => UserRole == "Admin" || UserRole == "Manager";
In your XAML:

XML
<Button Content="Add Item" IsEnabled="{Binding CanModifyData}" />
If an unauthorized user attempts a blocked action returned from the API with a 403 Forbidden status code, bind your view error banner to display: "Action Error: You lack sufficient permissions to perform this operation."Screen 1 - The Login Screen
The login window handles user credentials, calls the ApiService, manages validation/error states, and transitions to the main dashboard upon success.

Step 1: Create LoginViewModel.cs
In your ViewModels folder, create LoginViewModel.cs inheriting from ObservableObject:

C#
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using DesktopApp.Services;
using System.Threading.Tasks;

public partial class LoginViewModel : ObservableObject
{
    private readonly ApiService _apiService;

    [ObservableProperty]
    private string username;

    [ObservableProperty]
    private string password;

    [ObservableProperty]
    private string errorMessage;

    [ObservableProperty]
    private bool isErrorVisible;

    public IAsyncRelayCommand LoginCommand { get; }

    public LoginViewModel(ApiService apiService)
    {
        _apiService = apiService;
        LoginCommand = new AsyncRelayCommand(ExecuteLoginAsync);
    }

    private async Task ExecuteLoginAsync()
    {
        IsErrorVisible = false;
        
        if (string.IsNullOrWhiteSpace(Username) || string.IsNullOrWhiteSpace(Password))
        {
            ErrorMessage = "Please enter both username and password.";
            IsErrorVisible = true;
            return;
        }

        bool success = await _apiService.LoginAsync(Username, Password);
        if (success)
        {
            // Open Main Dashboard and close Login Window
            var mainWindow = new MainWindow(_apiService);
            mainWindow.Show();
            App.Current.Windows[0]?.Close();
        }
        else
        {
            ErrorMessage = "Invalid credentials or unauthorized access.";
            IsErrorVisible = true;
        }
    }
}
Step 2: Design LoginWindow.xaml
Update your primary startup window or create a dedicated LoginWindow.xaml:

XML
<Window x:Class="DesktopApp.Views.LoginWindow"
        xmlns="http://schemas.microsoft.com/winfx/2006/xaml/presentation"
        xmlns:x="http://schemas.microsoft.com/winfx/2006/xaml"
        Title="System Login" Height="350" Width="400" WindowStartupLocation="CenterScreen">
    <Grid Margin="20">
        <StackPanel VerticalAlignment="Center">
            <TextBlock Text="Welcome Back" FontSize="22" FontWeight="Bold" Margin="0,0,0,20"/>
            
            <TextBlock Text="Username" FontWeight="SemiBold"/>
            <TextBox Text="{Binding Username, UpdateSourceTrigger=PropertyChanged}" Height="30" Margin="0,5,0,15"/>
            
            <TextBlock Text="Password" FontWeight="SemiBold"/>
            <PasswordBox x:Name="PasswordBox" Height="30" Margin="0,5,0,15"/>

            <Button Content="Login" Command="{Binding LoginCommand}" Height="35" Background="#007ACC" Foreground="White" FontWeight="Bold"/>

            <TextBlock Text="{Binding ErrorMessage}" Foreground="Red" Margin="0,15,0,0" 
                       Visibility="{Binding IsErrorVisible, Converter={StaticResource BooleanToVisibilityConverter}}"/>
        </StackPanel>
    </Grid>
</Window>
Phase 4: Screen 2 & 3 - Dashboard with Inventory & Customer Management
The main dashboard (MainWindow.xaml) typically utilizes a sidebar or a TabControl to toggle between Inventory Management and Customer Management.

Step 1: Design MainWindow.xaml
Use a layout with navigation on the left or top tabs to switch contexts easily:

XML
<Window x:Class="DesktopApp.MainWindow"
        xmlns="http://schemas.microsoft.com/winfx/2006/xaml/presentation"
        xmlns:x="http://schemas.microsoft.com/winfx/2006/xaml"
        Title="Enterprise Management System" Height="600" Width="900" WindowStartupLocation="CenterScreen">
    <TabControl>
        <!-- Inventory Tab -->
        <TabItem Header="Inventory Management">
            <Grid Margin="10">
                <Grid.RowDefinitions>
                    <RowDefinition Height="Auto"/>
                    <RowDefinition Height="*"/>
                    <RowDefinition Height="Auto"/>
                </Grid.RowDefinitions>
                
                <StackPanel Orientation="Horizontal" Margin="0,0,0,10">
                    <Button Content="Add Item" Width="90" Height="30" Margin="0,0,10,0"/>
                    <Button Content="Edit Selected" Width="90" Height="30"/>
                </StackPanel>

                <DataGrid Grid.Row="1" AutoGenerateColumns="False" CanUserAddRows="False">
                    <DataGrid.Columns>
                        <DataGridTextColumn Header="Item ID" Width="100"/>
                        <DataGridTextColumn Header="Item Name" Width="250"/>
                        <DataGridTextColumn Header="Stock Quantity" Width="120"/>
                        <DataGridTextColumn Header="Price ($)" Width="100"/>
                    </DataGrid.Columns>
                </DataGrid>
                
                <TextBlock Grid.Row="2" Text="* Restricted permissions: Only Managers can edit inventory items." Foreground="Gray" Margin="0,10,0,0"/>
            </Grid>
        </TabItem>

        <!-- Customer Tab -->
        <TabItem Header="Customer Management">
            <Grid Margin="10">
                <Grid.RowDefinitions>
                    <RowDefinition Height="Auto"/>
                    <RowDefinition Height="*"/>
                </Grid.RowDefinitions>

                <StackPanel Orientation="Horizontal" Margin="0,0,0,10">
                    <Button Content="Add Customer" Width="110" Height="30" Margin="0,0,10,0"/>
                    <Button Content="Edit Customer" Width="110" Height="30"/>
                </StackPanel>

                <DataGrid Grid.Row="1" AutoGenerateColumns="False" CanUserAddRows="False">
                    <DataGrid.Columns>
                        <DataGridTextColumn Header="Customer ID" Width="100"/>
                        <DataGridTextColumn Header="Full Name" Width="250"/>
                        <DataGridTextColumn Header="Email" Width="200"/>
                        <DataGridTextColumn Header="Phone Number" Width="150"/>
                    </DataGrid.Columns>
                </DataGrid>
            </Grid>
        </TabItem>
    </TabControl>
</Window>
Phase 5: Handling Permissions & Action Errors
To handle role-based restrictions gracefully (e.g., Cashiers can view items, but only Admins can add/edit):

Decode the JWT token on login or store the Role returned in your LoginResponseDto.

Bind the IsEnabled property of your action buttons (Add/Edit) to a boolean property in your ViewModel (e.g., CanModifyData).

C#
// Example ViewModel logic snippet for permission checks
public bool CanModifyData => UserRole == "Admin" || UserRole == "Manager";
In your XAML:

XML
<Button Content="Add Item" IsEnabled="{Binding CanModifyData}" />
If an unauthorized user attempts a blocked action returned from the API with a 403 Forbidden status code, bind your view error banner to display: "Action Error: You lack sufficient permissions to perform this operation."