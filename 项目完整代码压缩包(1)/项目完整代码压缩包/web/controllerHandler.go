package web

import (
    "bytes"
    "crypto/md5"
    "database/sql"
    "encoding/base64"
    "encoding/json"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    shell "github.com/ipfs/go-ipfs-api"
    "html/template"
    "io"
    "io/ioutil"
    "math/rand"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "testwork/sdkInit"
    "time"
)


type FileData struct {
    Id                  string `json:"id"`
    Filename            string `json:"filename"`
    FileDetail          string `json:"filedetail"`
    FileHash            string `json:"fileHash"`
    CreateDate          string `json:"createDate"`
    Username            string `json:"username"`
    Card                string `json:"card"`
    User                []string `json:"user"`
    Filepath            string `json:"filepath"`
    IpfsHash            string `json:"ipfsHash"`
    Time                string `json:"time"`
    Date                string `json:"date"`
    Money               string `json:"money"`
    TransferDate        string `json:"transferdate"`
    Flag                string `json:"flag"`
    Tmp                string `json:"tmp"`
}
type Transaction struct {
    Id                  string `json:"id"`
    Filename            string `json:"filename"`
    Date                string `json:"date"`
    Username            string `json:"username"`
    User                string `json:"user"`
    Username1           string `json:"username1"`
    User1               string `json:"user1"`
    Money               string `json:"money"`
    FileHash               string `json:"fileHash"`
    Flag                string `json:"flag"`
}
type Login struct {
    name          string
    password      string
}
type Wallet struct {
    Id           string
    Address      string
    Money        string
}
type Ticket struct {
    Hash          string `json:"hash"`
    Money         string `json:"money"`
    Time          string `json:"time"`
    Content       string `json:"content"`
}
// //全局sdk变量 来操作链码
var App sdkInit.Application
//全局变量CA
var db *sql.DB
var sh *shell.Shell
//web服务路由信息
func UploadIPFS(str string) string {
    sh = shell.NewShell("localhost:5001")
    hash, err := sh.Add(bytes.NewBufferString(str))
    if err != nil {
        fmt.Println("ipfs错误：", err)
    }
    return hash
}

//从ipfs下载数据
func CatIPFS(hash string) string {
    sh = shell.NewShell("localhost:5001")
    read, err := sh.Cat(hash)
    if err != nil {
        fmt.Println(err)
    }
    body, err := ioutil.ReadAll(read)
    return string(body)
}
func initDB() (err error) {
    // DSN:Data Source Name
    dsn := "root:root123@tcp(localhost:3306)/Data"
    // 不会校验账号密码是否正确
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        return err
    }
    // 尝试与数据库建立连接（校验dsn是否正确）
    err = db.Ping()
    if err != nil {
        return err
    }
    return nil
}
//渲染模版
func ShowView(w http.ResponseWriter, r *http.Request, templateName string, data interface{})  {
    // 指定视图所在路径
    pagePath := filepath.Join("web/", templateName)
    resultTemplate, err := template.ParseFiles(pagePath)
    if err != nil {
        fmt.Printf("Error creating template instance: %v", err)
        return
    }
    //渲染模版
    err = resultTemplate.Execute(w, data)
    if err != nil {
        fmt.Printf("An error occurred while fusing data in the template: %v", err)
        return
    }
}
//将字符串进行md5加密
func MD5(str string) string {
    data := []byte(str) //切片
    has := md5.Sum(data)
    md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
    return md5str
}
func in(target string, str_array []string) bool {
    for _, element := range str_array{
        if target == element{
            return true
        }
    }
    return false
}
//admin
//登录函数
func a_login(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")
    index := r.FormValue("index")
    //输入为空 则返回原页面
    if username == "" && index == "" {
        ShowView(w, r, "admin/login.html", 0)
        //index也是返回原页面
    }else if index == "logout"{
        c := http.Cookie{
            Name: "name",
            Value: "",
        }
        //设置cookie
        http.SetCookie(w, &c)
        ShowView(w, r, "admin/login.html", 0)
    }else{
        //其他情形 判断是否账号密码一样
        var l Login
        var uid string
        err := initDB()
        if err != nil {
            fmt.Println("init failed, err:%v\n", err)
            return
        }
        defer db.Close()
        sql := "SELECT password,uid FROM user WHERE username = '" + username +"'"
        // fmt.Println(sql)
        //查询
        err = db.QueryRow(sql).Scan(&l.password, &uid)
        // fmt.Println(l.password)
        if err != nil {
            fmt.Println("login failed, err:%v\n", err)
        }
        //判断
        if password == l.password{
            c := http.Cookie{
                Name: "name",
                Value: uid,
            }
            //登录成功后 设置cookie
            http.SetCookie(w, &c)
            ShowView(w, r, "admin/login.html", 1)
        }else{
            ShowView(w, r, "admin/login.html", 2)
        }
    }
}

//注册
//用户进行注册
func a_register(w http.ResponseWriter, r *http.Request){
    name := r.FormValue("username")
    passwd := r.FormValue("password")
    uname := r.FormValue("uname")
    card := r.FormValue("card")
    fmt.Println(name, passwd, uname)
    rands := time.Now().Format("2006/01/02 15:04:05") + strconv.Itoa(rand.Intn(100))
    uid := "U-" + MD5(rands)
    sqlStr := "insert into user(uid, username, name, card, password) values (?,?,?,?,?)"
    err := initDB()
    if err != nil {
        fmt.Println("init failed, err:%v\n", err)
    }
    _, err = db.Exec(sqlStr, uid, name,uname,card, passwd)
    if err != nil {
        fmt.Errorf("insert failed, err:%v\n", err)
    }
    opera := []string{"setMoney", uid, "10", "init"}
    _, err = App.SetV(opera)
    if err != nil {
        fmt.Errorf("read failed, err:%v\n", err)
        fmt.Fprintln(w, "注册失败！")
        return
    }
    fmt.Fprintln(w, "注册成功！")
}

//获取cookie
//根据cookie判断权限 然后进行相应的跳转
func a_index(w http.ResponseWriter, r *http.Request) {
    name, err := r.Cookie("name")
    if err != nil {
        ShowView(w, r, "admin/403.html", nil)
        return
    } else if name.Value == ""{
        ShowView(w, r, "admin/403.html", nil)
    } else {
        ShowView(w, r, "admin/index.html", nil)
    }
}
func a_tickit(w http.ResponseWriter, r *http.Request){
    var ticket Ticket
    money := r.FormValue("money")
    ticket.Hash = MD5("ticket")
    ticket.Time = time.Now().Format("2006/01/02 15:04:05")
    ticket.Money = money
    putData, err := json.Marshal(ticket)
    if err != nil {
        fmt.Errorf("Failed to json asset: %s", err)
    }
    encodeData := base64.StdEncoding.EncodeToString(putData)
    ticket.Content = encodeData
    //fmt.Fprintln(w, "充值券: " + string(encodeData))
    ShowView(w, r, "admin/faq.html", ticket)
}

func a_money(w http.ResponseWriter, r *http.Request){
    var ticket Ticket
    cookie,_ :=r.Cookie("name")
    name := r.FormValue("money")
    decodeDataByteArr, err := base64.StdEncoding.DecodeString(name)
    if err != nil {
        fmt.Errorf("Failed to get asset: %s", err)
    }
    err = json.Unmarshal(decodeDataByteArr, &ticket)
    if ticket.Hash != MD5("ticket"){
        fmt.Errorf("setMoney failed, err:%v\n", err)
        fmt.Fprintln(w, "充值失败！")
        return
    }
    opera := []string{"setMoney", cookie.Value, ticket.Money, "up"}
    _, err = App.SetV(opera)
    if err != nil {
        fmt.Errorf("setMoney failed, err:%v\n", err)
        fmt.Fprintln(w, "充值失败！")
        return
    }
    fmt.Fprintln(w, "充值成功！")
}
func a_confirm(w http.ResponseWriter, r *http.Request) {
    name := r.FormValue("Name")
    var fd FileData
    opera := []string{"getFileDetail", name}
    value, err := App.Get(opera)
    if err != nil {
        fmt.Errorf("get failed, err:%v\n", err)
    }
    err = json.Unmarshal([]byte(value), &fd)
    if err != nil {
        fmt.Errorf("Failed to json asset: %s", err)
    }

    ShowView(w, r, "admin/provePrint.html", fd)
}

func a_delete(w http.ResponseWriter, r *http.Request) {
    name := r.FormValue("Name")
    cookie, _ := r.Cookie("name")
    err := initDB()
    if err != nil {
        fmt.Println("init failed, err:%v\n", err)
        return
    }
    defer db.Close()
    sql := "DELETE FROM transfer WHERE uid = '"+ name +"'"
    _, err = db.Exec(sql)
    if err != nil {
        fmt.Println("delet failed, err:%v\n", err)
    }
    opera := []string{"delete", name,""}
    _, err = App.Set(opera)
    if err != nil {
        fmt.Errorf("setMoney failed, err:%v\n", err)
        fmt.Fprintln(w, "失败！")
        return
    }
    var fileDatas []FileData
    var nfileDatas []FileData
    opera = []string{"getFileData", ""}
    value, err := App.Get(opera)
    if err != nil {
        fmt.Errorf("getFileData failed, err:%v\n", err)
        //return
    }
    err = json.Unmarshal([]byte(value), &fileDatas)
    if err != nil {
        fmt.Errorf("json failed, err:%v\n", err)
        //return
    }
    //not 未交易
    for _,fd := range fileDatas{
        if fd.Flag == "未交易"{
            continue
        }
        l := len(fd.User)
        if cookie.Value == fd.User[l-1] {
            nfileDatas = append(nfileDatas, fd)
        }
    }
    ShowView(w, r, "admin/projects.html", nfileDatas)
}
func a_route(w http.ResponseWriter, r *http.Request) {
    name := r.FormValue("Name")
    cookie, _ := r.Cookie("name")
    //name为list 则展示所有的项目
    if name == "lists" {
        var fileDatas []FileData
        var nfileDatas []FileData
        opera := []string{"getFileData", ""}
        value, err := App.Get(opera)
        if err != nil {
            fmt.Errorf("getFileData failed, err:%v\n", err)
            //return
        }
        err = json.Unmarshal([]byte(value), &fileDatas)
        if err != nil {
            fmt.Errorf("json failed, err:%v\n", err)
            //return
        }
        for _, fd := range fileDatas {
            l := len(fd.User)
            if cookie.Value == fd.User[l-1]{
                nfileDatas = append(nfileDatas, fd)
            }
        }
        ShowView(w, r, "admin/lists.html", nfileDatas)
    } else if name == "confirm"{
        var fileDatas []FileData
        var nfileDatas []FileData
        opera := []string{"getFileData", ""}
        value, err := App.Get(opera)
        if err != nil {
            fmt.Errorf("getFileData failed, err:%v\n", err)
            //return
        }
        err = json.Unmarshal([]byte(value), &fileDatas)
        if err != nil {
            fmt.Errorf("json failed, err:%v\n", err)
            //return
        }
        //show
        for _,fd := range fileDatas{
            l := len(fd.User)
            a := cookie.Value
            fmt.Println(a)
            if cookie.Value == fd.User[l-1]{
                if fd.Flag == "已交易"{
                    fd.Flag = "未交易"
                }
                nfileDatas = append(nfileDatas, fd)
            }else if in(cookie.Value, fd.User){
                fd.Flag = "已交易"
                nfileDatas = append(nfileDatas, fd)
            }
        }
        ShowView(w, r, "admin/confirm.html", nfileDatas)
    } else if name == "transaction"{
        var fileDatas []FileData
        var nfileDatas []FileData
        opera := []string{"getFileData", ""}
        value, err := App.Get(opera)
        if err != nil {
            fmt.Errorf("getFileData failed, err:%v\n", err)
            //return
        }
        err = json.Unmarshal([]byte(value), &fileDatas)
        if err != nil {
            fmt.Errorf("json failed, err:%v\n", err)
            //return
        }
        //show transaction
        for _,fd := range fileDatas{
            if len(fd.User) == 1{
                continue
            }
            l := len(fd.User)
            if cookie.Value == fd.User[l-1]{
                nfileDatas = append(nfileDatas, fd)
            }
        }
        ShowView(w, r, "admin/process.html", nfileDatas)
    } else if name == "operation"{
        var fileDatas []FileData
        var nfileDatas []FileData
        opera := []string{"getFileData", ""}
        value, err := App.Get(opera)
        if err != nil {
            fmt.Errorf("getFileData failed, err:%v\n", err)
            //return
        }
        err = json.Unmarshal([]byte(value), &fileDatas)
        if err != nil {
            fmt.Errorf("json failed, err:%v\n", err)
            //return
        }
        //not 未交易
        for _,fd := range fileDatas{
            if fd.Flag == "未交易"{
                continue
            }else if fd.Flag == "待交易"{
                l := len(fd.User)
                if cookie.Value == fd.User[l-1] {
                    nfileDatas = append(nfileDatas, fd)
                }
            }

        }
        ShowView(w, r, "admin/projects.html", nfileDatas)
    } else if name == "info"{
        var fileDatas []FileData
        var nfileDatas []FileData
        opera := []string{"getFileData", ""}
        value, err := App.Get(opera)
        if err != nil {
            fmt.Errorf("getFileData failed, err:%v\n", err)
            //return
        }
        err = json.Unmarshal([]byte(value), &fileDatas)
        if err != nil {
            fmt.Errorf("json failed, err:%v\n", err)
            //return
        }
        for _,fd := range fileDatas{
            if fd.Flag == "待交易" {
                nfileDatas = append(nfileDatas, fd)
            }
        }
        ShowView(w, r, "admin/project.html", nfileDatas)
    }else if name == "getMoney"{
        ShowView(w, r, "admin/getmoney.html", nil)
    } else if name == "add"{
        ShowView(w, r, "admin/add.html", nil)
    } else if name == "recharge"{
        var wallet Wallet
        opera := []string{"getMoney", cookie.Value, "-1"}
        value, _ := App.GetV(opera)
        wallet.Money = value
        wallet.Id = cookie.Value
        wallet.Address = MD5(wallet.Id)
        ShowView(w, r, "admin/wallet.html", wallet)
    }

}
func a_recharge(w http.ResponseWriter, r *http.Request) {
    ShowView(w, r, "admin/recharge.html", nil)
}

func a_apply(w http.ResponseWriter, r *http.Request) {
    name := r.FormValue("Name")
    money := r.FormValue("Money")
    var card string
    cookie, _ :=r.Cookie("name")
    date := time.Now().Format("2006/01/02 15:04:05")
    opera := []string{"getMoney", cookie.Value, money}
    value, err := App.GetV(opera)
    if err != nil {
        fmt.Errorf("get failed, err:%v\n", err)
        opera = []string{"getMoney", cookie.Value, "-1"}
        value, _ = App.GetV(opera)
        ShowView(w, r, "admin/found.html", value)
        return
    }
    err = initDB()
    if err != nil {
        fmt.Println("init failed, err:%v\n", err)
        return
    }
    sql := "SELECT name,card FROM user WHERE uid = '" + cookie.Value +"'"
    // fmt.Println(sql)
    //查询
    var username string
    err = db.QueryRow(sql).Scan(&username,&card)
    // fmt.Println(l.password)
    if err != nil {
        fmt.Println("login failed, err:%v\n", err)
    }

    opera = []string{"setTransaction", cookie.Value, name, username, date, card}
    value, err = App.SetVV(opera)
    if err != nil {
        fmt.Errorf("get failed, err:%v\n", err)
    }
    sql = "DELETE FROM transfer WHERE uid = '"+ name +"'"
    _, err = db.Exec(sql)
    if err != nil {
        fmt.Println("delet failed, err:%v\n", err)
    }
    ShowView(w, r, "admin/success.html", value)
}
func a_detail(w http.ResponseWriter, r *http.Request) {
    name := r.FormValue("Name")
    set := r.FormValue("Set")
    fmt.Println(name,set)
    var fd FileData
    //获取用户名字
    var ts []Transaction
    if set == "transfer"{
        opera := []string{"getTransaction", name}
        value, err := App.Get(opera)
        if err != nil {
            fmt.Errorf("get failed, err:%v\n", err)
        }
        err = json.Unmarshal([]byte(value), &ts)
        if err != nil {
            fmt.Errorf("Failed to json asset: %s", err)
        }
        ShowView(w, r, "admin/transaction.html", ts)
    }else {
        opera := []string{"getFileDetail", name}
        value, err := App.Get(opera)
        if err != nil {
            fmt.Errorf("get failed, err:%v\n", err)
        }
        err = json.Unmarshal([]byte(value), &fd)
        if err != nil {
            fmt.Errorf("Failed to json asset: %s", err)
        }
        fd.Flag = set
        ShowView(w, r, "admin/cat.html", fd)
    }
}
func a_download(rw http.ResponseWriter,r *http.Request){
    //获取请求参数
    fn :=r.FormValue("filename")
    fp :=r.FormValue("filepath")
    data := CatIPFS(fn) //return []byte 文件字节流
    //设置响应头
    header:=rw.Header()
    header.Add("Content-Type","application/octet-stream")
    header.Add("Content-Disposition","attachment;filename="+fp)
    //写入到响应流中
    rw.Write([]byte(data))
}
//项目增加
func a_add(w http.ResponseWriter, r *http.Request) {
    cookie, _:=r.Cookie("name")
    var fileData FileData
    ///获取其他值
    name := r.FormValue("name")
    detail := r.FormValue("detail")
    //文件处理
    r.ParseMultipartForm(32 << 20)
    file, handler, err := r.FormFile("file")
    filepath := handler.Filename
    fmt.Println(filepath)
    if err != nil {
        fmt.Fprintln(w, "上传失败！")
        fmt.Println(err)
        return
    }
    defer file.Close()
    fileContent, _ :=ioutil.ReadAll(file)
    fileData.FileHash = UploadIPFS(string(fileContent))
    times := strconv.FormatInt(time.Now().Unix(), 10)

    //ipfs操作
    fmt.Println("文件hash是", fileData.FileHash)
    //ipfs索引hash加密
    f, err := os.OpenFile("./web/static/upload/"+ MD5(fileData.FileHash), os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        fmt.Fprintln(w, "update error！")
        fmt.Println(err)
        return
    }
    defer f.Close()
    io.Copy(f, file)
    //形成fileData
    fileData.CreateDate = time.Now().Format("2006/01/02 15:04:05")
    fileData.Date = time.Now().Format("2006/01/02 15:04:05")
    rands := time.Now().Format("2006/01/02 15:04:05") + strconv.Itoa(rand.Intn(100))
    fileData.Id = MD5(rands)
    fileData.Filename = name
    fileData.Card = r.FormValue("card")
    fileData.FileDetail = detail
    fileData.Filepath = filepath
    fileData.User = []string{cookie.Value}
    fileData.Time = times
    fileData.Flag = "未交易"
    fileData.Tmp = "未交易"

    r.ParseMultipartForm(32 << 20)
    file, handler, err = r.FormFile("files")
    filepath = handler.Filename
    fmt.Println(filepath)
    if err != nil {
        fmt.Fprintln(w, "上传失败！")
        fmt.Println(err)
        return
    }
    defer file.Close()
    fileContent, _ =ioutil.ReadAll(file)
    if err := ioutil.WriteFile("./web/static/img/"+fileData.Id, fileContent, 0666); err != nil {
        fmt.Fprintln(w, "错误！")
        fmt.Println(err)
        return
    }
    
    err = initDB()
    if err != nil {
        fmt.Println("init failed, err: %v\n", err)
        return
    }
    sql := "SELECT name FROM user WHERE uid = '" + cookie.Value +"'"
    err = db.QueryRow(sql).Scan(&fileData.Username)
    if err != nil {
        fmt.Errorf("select failed, err: %v\n", err)
    }
    //序列化
    putData, err := json.Marshal(fileData)
    if err != nil {
        fmt.Errorf("Failed to json asset: %s", err)
        fmt.Fprintln(w, "上传失败！")
        return
    }
    //上传区块链文件
    opera := []string{"setFileData", string(putData), ""}
    _, err = App.Set(opera)
    if err != nil {
        fmt.Errorf("setFileData failed, err:%v\n", err)
        fmt.Fprintln(w, "上传失败！")
        return
    }
    fmt.Fprintln(w, "上传成功！")
}
