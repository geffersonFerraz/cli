package beautiful

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Output fornece funções para embelezar diferentes tipos de output
type Output struct {
	rawMode bool
}

// NewOutput cria uma nova instância do embelezador de output
func NewOutput(rawMode bool) *Output {
	return &Output{
		rawMode: rawMode,
	}
}

// PrintJSON embelezar output JSON com cores e formatação
func (bo *Output) PrintJSON(data interface{}) error {
	if bo.rawMode {
		// Modo raw: output simples sem formatação
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
		return nil
	}

	// Modo embelezado: JSON formatado com cores
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Aplicar cores ao JSON
	coloredJSON := bo.colorizeJSON(string(jsonData))
	fmt.Println(coloredJSON)
	return nil
}

// PrintSuccess embelezar mensagens de sucesso
func (bo *Output) PrintSuccess(message string) {
	if bo.rawMode {
		fmt.Println(message)
		return
	}

	successColor := color.New(color.FgGreen, color.Bold)
	successColor.Printf("✅ %s\n", message)
}

// PrintError embelezar mensagens de erro
func (bo *Output) PrintError(message string) {
	if bo.rawMode {
		fmt.Fprintf(os.Stderr, "Error: %s\n", message)
		return
	}

	errorColor := color.New(color.FgRed, color.Bold)
	errorColor.Printf("❌ %s\n", message)
}

// PrintWarning embelezar mensagens de aviso
func (bo *Output) PrintWarning(message string) {
	if bo.rawMode {
		fmt.Printf("Warning: %s\n", message)
		return
	}

	warningColor := color.New(color.FgYellow, color.Bold)
	warningColor.Printf("⚠️  %s\n", message)
}

// PrintInfo embelezar mensagens informativas
func (bo *Output) PrintInfo(message string) {
	if bo.rawMode {
		fmt.Println(message)
		return
	}

	infoColor := color.New(color.FgCyan, color.Bold)
	infoColor.Printf("ℹ️  %s\n", message)
}

// PrintTable embelezar dados em formato de tabela
func (bo *Output) PrintTable(headers []string, rows [][]string) {
	if bo.rawMode {
		// Modo raw: output simples
		fmt.Println(strings.Join(headers, "\t"))
		for _, row := range rows {
			fmt.Println(strings.Join(row, "\t"))
		}
		return
	}

	// Modo embelezado: tabela com cores e bordas
	headerColor := color.New(color.FgMagenta, color.Bold)
	rowColor := color.New(color.FgWhite)

	// Calcular larguras das colunas
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Imprimir cabeçalho
	headerColor.Print("┌")
	for i, width := range widths {
		for j := 0; j < width+2; j++ {
			headerColor.Print("─")
		}
		if i < len(widths)-1 {
			headerColor.Print("┬")
		}
	}
	headerColor.Println("┐")

	// Imprimir títulos das colunas
	headerColor.Print("│")
	for i, header := range headers {
		headerColor.Printf(" %-*s │", widths[i], header)
	}
	headerColor.Println()

	// Imprimir separador
	headerColor.Print("├")
	for i, width := range widths {
		for j := 0; j < width+2; j++ {
			headerColor.Print("─")
		}
		if i < len(widths)-1 {
			headerColor.Print("┼")
		}
	}
	headerColor.Println("┤")

	// Imprimir linhas de dados
	for _, row := range rows {
		rowColor.Print("│")
		for i, cell := range row {
			if i < len(widths) {
				rowColor.Printf(" %-*s │", widths[i], cell)
			}
		}
		rowColor.Println()
	}

	// Imprimir rodapé
	headerColor.Print("└")
	for i, width := range widths {
		for j := 0; j < width+2; j++ {
			headerColor.Print("─")
		}
		if i < len(widths)-1 {
			headerColor.Print("┴")
		}
	}
	headerColor.Println("┘")
}

// PrintList embelezar listas
func (bo *Output) PrintList(title string, items []string) {
	if bo.rawMode {
		fmt.Println(title)
		for _, item := range items {
			fmt.Printf("- %s\n", item)
		}
		return
	}

	titleColor := color.New(color.FgBlue, color.Bold)
	titleColor.Printf("📋 %s:\n", title)

	itemColor := color.New(color.FgCyan)
	for i, item := range items {
		itemColor.Printf("  %d. %s\n", i+1, item)
	}
}

// colorizeJSON aplica cores ao JSON formatado
func (bo *Output) colorizeJSON(jsonStr string) string {
	if bo.rawMode {
		return jsonStr
	}

	// Cores para diferentes elementos do JSON
	braceColor := color.New(color.FgWhite, color.Bold)
	keyColor := color.New(color.FgYellow, color.Bold)
	stringColor := color.New(color.FgGreen)
	numberColor := color.New(color.FgCyan)
	booleanColor := color.New(color.FgMagenta, color.Bold)
	nullColor := color.New(color.FgRed, color.Bold)

	lines := strings.Split(jsonStr, "\n")
	var result []string

	for _, line := range lines {
		// Aplicar cores baseadas no conteúdo da linha
		coloredLine := bo.colorizeJSONLine(line, braceColor, keyColor, stringColor, numberColor, booleanColor, nullColor)
		result = append(result, coloredLine)
	}

	return strings.Join(result, "\n")
}

// colorizeJSONLine aplica cores a uma linha específica do JSON
func (bo *Output) colorizeJSONLine(line string, braceColor, keyColor, stringColor, numberColor, booleanColor, nullColor *color.Color) string {
	if bo.rawMode {
		return line
	}

	// Identificar e colorir diferentes elementos
	line = strings.ReplaceAll(line, `"`, stringColor.Sprint(`"`))
	line = strings.ReplaceAll(line, `{`, braceColor.Sprint(`{`))
	line = strings.ReplaceAll(line, `}`, braceColor.Sprint(`}`))
	line = strings.ReplaceAll(line, `[`, braceColor.Sprint(`[`))
	line = strings.ReplaceAll(line, `]`, braceColor.Sprint(`]`))
	line = strings.ReplaceAll(line, `,`, braceColor.Sprint(`,`))
	line = strings.ReplaceAll(line, `:`, braceColor.Sprint(`:`))

	// Colorir valores booleanos
	line = strings.ReplaceAll(line, `true`, booleanColor.Sprint(`true`))
	line = strings.ReplaceAll(line, `false`, booleanColor.Sprint(`false`))
	line = strings.ReplaceAll(line, `null`, nullColor.Sprint(`null`))

	return line
}

// PrintProgress embelezar barras de progresso
func (bo *Output) PrintProgress(current, total int, message string) {
	if bo.rawMode {
		fmt.Printf("%s: %d/%d\n", message, current, total)
		return
	}

	progressColor := color.New(color.FgBlue, color.Bold)
	progressColor.Printf("🔄 %s: %d/%d\n", message, current, total)
}

// PrintHeader embelezar cabeçalhos de seção
func (bo *Output) PrintHeader(title string) {
	if bo.rawMode {
		fmt.Printf("\n=== %s ===\n", title)
		return
	}

	headerColor := color.New(color.FgCyan, color.Bold)
	headerColor.Printf("\n🎯 %s\n", title)
	headerColor.Println(strings.Repeat("─", len(title)+4))
}
