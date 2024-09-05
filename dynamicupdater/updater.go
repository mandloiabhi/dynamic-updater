package dynamicupdater

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	pluginDir       = "./dynamicupdater/plugins" // Path to the plugins directory
	loadedFunctions = make(map[string]interface{})
	functionSignatures = make(map[string]reflect.Type) // Store expected signatures dynamically
	mu              sync.Mutex
)

// Initialize sets up the plugin monitoring routine and creates the necessary directories
func Initialize() {
	// Create the plugin directory if it doesn't exist
	fmt.Println("ok Initialize function is working")
	if err := os.MkdirAll(pluginDir, os.ModePerm); err != nil {
		fmt.Println("Error creating plugin directory:", err)
		return
	}

	// Start a goroutine to monitor plugins
	go monitorPlugins()

	fmt.Println("Dynamic updater initialized. Monitoring for plugin changes...")
}

// monitorPlugins monitors the plugin directory for changes
func monitorPlugins() {
	for {
		mu.Lock()
		loadPlugins()
		mu.Unlock()
		time.Sleep(2 * time.Second) // Adjust the polling interval as needed
	}
}

// loadPlugins loads all plugins from the plugins directory
func loadPlugins() {
	files, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		fmt.Println("Error reading plugins directory:", err)
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".so" {
			pluginPath := filepath.Join(pluginDir, file.Name())
			p, err := plugin.Open(pluginPath)
			if err != nil {
				fmt.Println("Error loading plugin:", err)
				continue
			}

			functionName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			functionName = strings.TrimSuffix(functionName, ".go")
			fmt.println(functionName)
			f, err := p.Lookup(functionName)
			if err != nil {
				fmt.Println("Error finding function in plugin:", err)
				continue
			}

			// Use the reflect package to dynamically validate the function's signature
			if validateFunctionSignature(functionName, f) {
				loadedFunctions[functionName] = f
				functionSignatures[functionName] = reflect.TypeOf(f) // Store the function signature dynamically
				fmt.Printf("Loaded function: %s from plugin: %s\n", functionName, file.Name())
			} else {
				fmt.Printf("Invalid signature for function: %s\n", functionName)
			}
		}
	}
}

// GetFunction retrieves a dynamically loaded function by name
func GetFunction(name string) (interface{}, error) {
	mu.Lock()
	defer mu.Unlock()

	if fn, ok := loadedFunctions[name]; ok {
		return fn, nil
	}
	return nil, fmt.Errorf("function %s not found", name)
}

// calculateFunctionSignature computes the SHA-256 hash of the function code (for versioning or change detection)
func calculateFunctionSignature(pluginPath string) string {
	data, err := ioutil.ReadFile(pluginPath)
	if err != nil {
		fmt.Println("Error reading plugin file:", err)
		return ""
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// validateFunctionSignature validates the dynamically loaded function signature
func validateFunctionSignature(functionName string, loadedFunction interface{}) bool {
	mu.Lock()
	defer mu.Unlock()

	// Get the reflect.Type of the loaded function
	loadedFuncType := reflect.TypeOf(loadedFunction)

	// Check if we have a stored signature for this function
	if expectedType, exists := functionSignatures[functionName]; exists {
		// Compare the function signatures dynamically
		return loadedFuncType == expectedType
	}

	// If no signature exists, store it and allow this as the first valid instance
	functionSignatures[functionName] = loadedFuncType
	return true
}


// package dynamicupdater

// import (
// 	"crypto/sha256"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"plugin"
// 	"strings"
// 	"sync"
// 	"time"
// )

// var (
// 	pluginDir          = "./dynamicupdater/plugins" // Path to the plugins directory
// 	loadedFunctions    = make(map[string]interface{})
// 	functionSignatures = make(map[string]string) // Store signatures for validation
// 	mu                 sync.Mutex
// )

// // Initialize sets up the plugin monitoring routine and creates the necessary directories
// func Initialize() {
// 	// Create the plugin directory if it doesn't exist
// 	fmt.Println("ok Initialize function is working")
// 	if err := os.MkdirAll(pluginDir, os.ModePerm); err != nil {
// 		fmt.Println("Error creating plugin directory:", err)
// 		return
// 	}

// 	// Start a goroutine to monitor plugins
// 	go monitorPlugins()

// 	fmt.Println("Dynamic updater initialized. Monitoring for plugin changes...")
// }

// // monitorPlugins monitors the plugin directory for changes
// func monitorPlugins() {
// 	for {
// 		mu.Lock()
// 		loadPlugins()
// 		mu.Unlock()
// 		time.Sleep(5 * time.Second) // Adjust the polling interval as needed
// 	}
// }

// // loadPlugins loads all plugins from the plugins directory
// func loadPlugins() {
// 	files, err := ioutil.ReadDir(pluginDir)
// 	if err != nil {
// 		fmt.Println("Error reading plugins directory:", err)
// 		return
// 	}

// 	for _, file := range files {
// 		if filepath.Ext(file.Name()) == ".so" {
// 			pluginPath := filepath.Join(pluginDir, file.Name())
// 			p, err := plugin.Open(pluginPath)
// 			if err != nil {
// 				fmt.Println("Error loading plugin:", err)
// 				continue
// 			}

// 			functionName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
// 			functionName = strings.TrimSuffix(functionName, ".go")
// 			f, err := p.Lookup(functionName)
// 			if err != nil {
// 				fmt.Println("Error finding function in plugin:", err)
// 				continue
// 			}

// 			if validateFunctionSignature(pluginPath, functionName) {
// 				loadedFunctions[functionName] = f
// 				functionSignatures[functionName] = calculateFunctionSignature(pluginPath, functionName) // Save the signature for validation
// 				fmt.Printf("Loaded function: %s from plugin: %s\n", functionName, file.Name())
// 			} else {
// 				fmt.Printf("Invalid signature for function: %s\n", functionName)
// 			}
// 		}
// 	}
// }

// // GetFunction retrieves a dynamically loaded function by name
// func GetFunction(name string) (interface{}, error) {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	if fn, ok := loadedFunctions[name]; ok {
// 		return fn, nil
// 	}
// 	return nil, fmt.Errorf("function %s not found", name)
// }

// // calculateFunctionSignature computes the SHA-256 hash of the function code
// func calculateFunctionSignature(pluginPath string, functionName string) string {
// 	data, err := ioutil.ReadFile(pluginPath)
// 	if err != nil {
// 		fmt.Println("Error reading plugin file:", err)
// 		return ""
// 	}
// 	hash := sha256.Sum256(data)
// 	return fmt.Sprintf("%x", hash)
// }

// // validateFunctionSignature checks the signature of the function against the expected value
// func validateFunctionSignature(pluginPath string, functionName string) bool {
// 	expectedSignature := functionSignatures[functionName]
// 	currentSignature := calculateFunctionSignature(pluginPath, functionName)
// 	return expectedSignature == currentSignature
// }
