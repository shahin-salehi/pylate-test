// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.898
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "shahin/webserver/internal/web/components"

func Login() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.Head("Login | Shahin").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<body class=\"min-h-screen flex items-center justify-center bg-gradient-to-br from-green-50 to-green-100\"><div class=\"w-full max-w-md bg-white shadow-2xl p-6 rounded-2xl\"><h1 class=\"text-3xl font-bold text-center mb-6\">Welcome</h1><form method=\"POST\" action=\"/login\" class=\"space-y-4\"><div><label for=\"email\" class=\"block text-sm font-medium\">Email</label> <input type=\"email\" id=\"email\" name=\"email\" placeholder=\"you@example.com\" class=\"mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-green-500 focus:ring-green-500\"></div><div><label for=\"password\" class=\"block text-sm font-medium\">Password</label> <input type=\"password\" id=\"password\" name=\"password\" placeholder=\"••••••••\" class=\"mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-green-500 focus:ring-green-500\"></div><button type=\"submit\" class=\"w-full px-6 py-3 text-white rounded-xl transition rounded-md\" style=\"background-color: #003824;\" onmouseover=\"this.style.backgroundColor='#002b1c'\" onmouseout=\"this.style.backgroundColor='#003824'\">Sign In</button><!--class=\"w-full bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded-md\"--></form><p class=\"text-sm text-center text-gray-600 mt-4\">Don’t have an account? Contact <a href=\"/signup\" class=\"text-green-600 font-medium hover:underline\">AIOps</a></p></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.Footer("/static/js/login/script.js").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
