package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"net/http"

	"path/filepath"

	"io/ioutil"

	"strconv"

	"log"

	"github.com/spf13/cobra"
	"gopkg.in/cheggaaa/pb.v1"
)

const (
	NAME       = "izip"
	URL        = "http://local.izip.softped.com"
	URL_ID     = URL + "/id"
	URL_UPLOAD = URL + "/upload"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   NAME + " <file>",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("No file specified.")
		}
		if len(args) > 1 {
			return fmt.Errorf("Only one file can be given.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "File '%s' doesn't exist.", args[0])
			return
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		if fi.IsDir() {
			fmt.Fprintf(os.Stderr, "'%s' is a directory.", args[0])
			return
		}

		req, err := http.NewRequest("POST", URL_ID, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		req.Header.Set("X-FileName", filepath.Base(f.Name()))
		req.Header.Set("X-FileSize", strconv.FormatInt(fi.Size(), 10))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			fmt.Fprintln(os.Stderr, resp.Status)
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Fprintf(os.Stderr, "%s", b)
			}
			return
		}

		var r idResponse
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return
		}
		fmt.Println(r.URL)

		bar := pb.New(int(fi.Size())).SetUnits(pb.U_BYTES)
		bar.ShowTimeLeft = true
		bar.ShowSpeed = true
		bar.Start()
		defer bar.Finish()

		req, err = http.NewRequest("POST", URL_UPLOAD, bar.NewProxyReader(f))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		req.Header.Set("X-ID", r.ID)

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			fmt.Fprintln(os.Stderr, resp.Status)
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Fprintf(os.Stderr, "%s", b)
			}
			return
		}
	},
}

type idResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	helpFunc := RootCmd.HelpFunc()
	RootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		initVerbose()
		CheckUpdate()
		helpFunc(cmd, args)
	})
	if _, _, err := RootCmd.Find(os.Args[1:]); err != nil {
		// append '--' in front of the args
		os.Args = append(os.Args, "")
		copy(os.Args[2:], os.Args[1:])
		os.Args[1] = "--"
	}
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

var (
	verbose bool
)

func init() {
	cobra.OnInitialize(initVerbose)
	cobra.OnInitialize(CheckUpdate)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "print logs")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initVerbose() {
	if !verbose {
		log.SetOutput(ioutil.Discard)
	}
}
