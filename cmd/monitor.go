package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Feinq/mleak/pkg/monitor"
	"github.com/Feinq/mleak/pkg/report"
	"github.com/spf13/cobra"
)

var (
	interval  time.Duration
	outputFmt string
)

var monitorCmd = &cobra.Command{
	Use:   "monitor [pid]",
	Short: "Monitor a process memory usage",
	Long:  `Continuously monitor the memory usage of a specified process identified by its PID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid PID: %s", args[0])
		}

		processes, err := monitor.GetProcesses()
		if err != nil {
			return err
		}

		processName := ""
		for _, p := range processes {
			if p.PID == pid {
				processName = p.Name
				break
			}
		}

		table := report.NewTableMemoryReport(pid, processName, interval, 10)
		trend := report.NewMemoryTrend(20)

		for {
			mem, err := monitor.GetProcessMemory(pid)
			if err != nil {
				return fmt.Errorf("failed to get memory for PID %d: %w", pid, err)
			}

			rss := float64(mem)
			now := time.Now()
			timestamp := now.Format(time.RFC3339)

			trend.AddSample(rss, now)

			risk := report.DetermineLeakRiskWithTrend(trend)

			r := report.MemoryReport{
				Timestamp: timestamp,
				RSS:       rss,
				LeakRisk:  risk,
			}

			if outputFmt == "json" {
				r.OutputJSON()
				fmt.Println()
			} else if outputFmt == "table" {
				table.AddReport(r)
				table.DisplayTable()
			} else if outputFmt == "text" {
				r.OutputPlainText()
			}

			time.Sleep(interval)
		}
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)

	monitorCmd.Flags().DurationVarP(&interval, "interval", "i", 10*time.Second, "Monitoring interval")
	monitorCmd.Flags().StringVarP(&outputFmt, "format", "f", "table", "Output format (table, text, json)")
}
