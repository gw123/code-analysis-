package usecase

import "strings"

func buildFileAnalysisPrompt(filename, code string) string {
	p := `请分析以下的代码文件，并提取相关信息。请注意以下要点：
1. **功能描述**
   - 总结代码文件的整体功能和用途，并列出所有可以导出的结构体、常量、接口的名称。
   
2. **文件基本信息**
   - 文件名：
   - 包名：
   - 依赖导入项目（列出所有导入的包）：

3. **常量**
   - 列出所有常量及其值，并简要描述功能。

4. **结构体**
   - 列出所有结构体及其字段与类型。
   - 列出每个结构体的所有方法（函数），并简要描述功能。

5. **Golang接口**
   - 列出所有接口及其方法，并简要描述每个方法的功能、参数和返回值。

6. **方法**
   - 列出所有方法及其参数和返回值。
   - 简要描述每个方法的功能。

7. **API接口(如果存在)**
   - 列出接口的请求参数。
   - 列出接口的响应格式。
   - 列出接口的请求方式: GET | POST | PUT | DELETE。

请逐项回答，确保信息清晰明了：

- 输出格式使用**YAML**结构化。
- 参考下面的输出格式：
- 保证输出内容只包含YAML结构，方便后续解析。
- 输出的描述信息使用中文。
- 对应字段的值如有混淆，使用单引号包裹。
- 确保格式清晰正确，保持与以下示例一致，便于代码解析。
- 若某些部分（如structs、constants、interfaces等）为空，不要输出对应字段。

**注意：**为便于理解，代码中会使用以下术语：
- **image:** 镜像
- **artifactory:** 制品仓库
- **artifact:** 制品

---

### 输出示例：
file_description: |
    <文件的功能是实现XXX>

    file_info:
	file_name: <file_name>
	package_name: <package_name>
	imports:
	- <package_1>
	- <package_2>

	constants:
	- name: <constant_name>
	value: <constant_value>
	description: <constant_function_description>

	structs:
	- name: <struct_name>
	fields:
	- '<field_1>: <type_1>'
	- '<field_2>: <type_2>'
	methods:
	- name: <method_name>
	params:
	- <param_1>
	return_values:
	- <return_type>
	description: <method_description>

	interfaces:
	- name: <interface_name>
	methods:
	- name: <method_name>
	params:
	- <param_1>
	return_values:
	- <return_type>
	description: <method_description>

	methods:
	- name: <method_name>
	params:
	- <param_1>
	return_values:
	- <return_type>
	description: <method_description>

	api_endpoints:
	- name: <api_name>
	request_params:
	- <param_1>
	response:
	- <response_format>
	request_method: '<GET|POST|PUT|DELETE>'
`

	strBuilder := strings.Builder{}
	strBuilder.WriteString(p)
	strBuilder.WriteString("文件名: ")
	strBuilder.WriteString(filename)
	strBuilder.WriteString("\n")
	strBuilder.WriteString("以下是代码文件：\n")
	strBuilder.WriteString(code)
	return strBuilder.String()
}

func buildQuestionRelFilesPrompt(question, summary string) string {
	strBuilder := strings.Builder{}
	strBuilder.WriteString(`你的角色是一个高级开发工程师。根据以下 Golang 源代码中各个文件的总结信息，请回答下面问题。`)
	strBuilder.WriteString(question)
	strBuilder.WriteString(`
	输出结果要求:
	1.只需要列出与该功能相关的文件和选择该文件的依据。
    2.请按照方法的调用层级从低到高输出
    3.只输出yaml内容

### 输出示例:
- file: '<xxx.go>'
  why: '<解释一下为啥选择这个文件>'

### 以下是源码信息:
	`)

	strBuilder.WriteString(summary)
	return strBuilder.String()
}

func buildQuestionRelFilesParsePrompt(question, step1Answer, filename, fileContent string) string {

	strBuilder := strings.Builder{}
	{
		strBuilder.WriteString(`你的角色是一个高级开发工程师。根据以下 Golang 源代码中相关文件的总结信息，回答下面问题:`)
		strBuilder.WriteString(question)
		strBuilder.WriteString(`
	### 输出结果要求:
    1.解释该源码中关键的方法和方法的作用
    2.列出方法的调用关系
    下面是第一步分析得到的总结信息:`)
		strBuilder.WriteString(step1Answer)
		strBuilder.WriteString("\n\n")
		strBuilder.WriteString(`
	### 以下是 ` + filename + `文件源码信息：
	`)

		strBuilder.WriteString(string(fileContent))
	}
	return strBuilder.String()
}

func buildFinalAnswerPrompt(question, helpInfo string) *strings.Builder {
	strBuilder3 := strings.Builder{}
	strBuilder3.WriteString(`你的角色是一个高级开发工程师。根据以下 Golang 源代码中相关文件的总结信息，回答下面问题:`)
	strBuilder3.WriteString(question)
	strBuilder3.WriteString(`
	### 输出结果要求:
    1. 输出一个 remind 图表示方法之间的调用关系
       输出示例:
       <xx>
		  |
		  v
		<xx>.go
		  |
		  v
		<xx>.go
		  |
		  +---> <xx>.go
		  |         |
		  |         +---> <xx>.go <xxfunction> <return_type>
		  |
		  +---> <xx>.go
    2. 总结功能实现的逻辑
    3. 如果问题中是需要实现一个功能,请写出实现的代码逻辑,以及代码放在什么地方合适
`)
	strBuilder3.WriteString(`
	### 以下是相关参考信息:
	`)
	strBuilder3.WriteString(helpInfo)
	strBuilder3.WriteString(`
	### 以下是文件源码信息：
	`)
	return &strBuilder3
}

func GenCodeUseDocHelpInfo() string {
	strBuilder := strings.Builder{}

	strBuilder.WriteString(`### 输出要求:
请根据源代码生成 auth 节点的使用说明文档，结构如下：

节点名称：auth
描述该节点的用途及在工作流中的角色。

节点类型 (type)：trigger | node

凭证信息 (credentials)：
说明节点需要哪些认证参数，如 api_key、smtp_host、username、password 等。确保每个节点都正确配置了凭证。

节点参数 (parameters)：
定义节点运行所需的参数，根据节点类型变化。例如，mail 节点的参数可能包括 subject、body、to。

节点输出：
节点的 Run 方法返回值。

示例：邮件节点配置
nodes:
  - name: sendMail
    type: mail
    credentials:
      smtp_host: "{{SMTP_HOST}}"
      smtp_port: {{SMTP_PORT}}
      username: "{{USERNAME}}"
      password: "{{PASSWORD}}"
      from: "{{FROM_EMAIL}}"
    parameters:
      subject: "{{SUBJECT}}"
      body: "{{BODY}}"
      to:
        - "{{TO_EMAIL1}}"
        - "{{TO_EMAIL2}}"
说明：
节点名称：sendMail，表示发送邮件的节点。
类型：mail，通过 SMTP 发送电子邮件。
凭证信息：配置 SMTP 服务器信息。
参数：邮件的主题、正文及收件人列表。


输出示例：
output:
   status: <表示成功执行>  
   data: <执行结果>
   error: <错误详情>
说明：
 status: xxx
 data: xxx
 error: xxx

`)
	return strBuilder.String()
}

func GenNodeHelpInfo() string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(`### 参考信息
节点接口定义
type NodeRunner interface {
    Init(node *entity.Node) error        // 初始化节点
    GetNode() *entity.Node               // 获取节点
    GetKind() string                     // 获取节点执行器名称
    GetType() string                     // 获取节点类型：Trigger | Node
    Run(ctx context.Context, entity *entity.WorkflowEntity) (map[string]interface{}, int, error) // 执行节点
    GetNodeParameter(entity *entity.WorkflowEntity, path, defaultValue string) string // 获取节点参数
}
父类提供的额外方法
func (b *BaseNode) GetNodeParameterInt(entity *entity.WorkflowEntity, path string, defaultValue int) int
func (b *BaseNode) GetNodeParameterFloat(entity *entity.WorkflowEntity, path string, defaultValue float64) float64
func (b *BaseNode) GetNodeParameter(entity *entity.WorkflowEntity, path, defaultValue string) string
func (b *BaseNode) GetNodeParameterArrayMap(entity *entity.WorkflowEntity, path string, defaultValue []map[string]interface{}) []map[string]interface{}
func (b *BaseNode) GetNodeParameterMap(entity *entity.WorkflowEntity, path string, defaultValue map[string]string) map[string]string
func (b *BaseNode) GetNodeParameterArray(entity *entity.WorkflowEntity, path string, defaultValue []string) []string
func (b *BaseNode) GenerateArtifact(ctx context.Context, workflowEntity *entity.WorkflowEntity, body io.ReadCloser) (*storage.UploadResult, error)
func (b *BaseNode) EvalStrTpl(entity *entity.WorkflowEntity, strTpl, defaultValue string) string

初始化示例
定义初始化参数结构体。
在 Init 函数中解析节点初始化参数。例如 Redis 节点：
type RedisOptions struct {
    Type     string
    Addr     string
    Password string
    DB       int
}

func (r *Redis) Init(node *entity.Node) error {
    err := r.BaseNode.Init(node)
    if err != nil {
        return err
    }
    err = json.Unmarshal(node.Credentials, &r.options)
    if err != nil {
        return app_err.NewNodeInitErr("redis", "BaseNode Init failed", err)
    }

    r.client = redis.NewClient(&redis.Options{
        Addr:     r.options.Addr,
        Password: r.options.Password,
        DB:       r.options.DB,
    })
    return nil
}

节点返回执行结果给下一个节点
func (e *EmailNode) Run(ctx context.Context, entity *entity.WorkflowEntity) (map[string]interface{}, int, error) {
    ....
    // 返回的第一个参数就是下一个节点的输入也就是 {{$last.xxx}} 中的 $last
	return map[string]interface{}{"status": "success"}, 0, nil
}


`)
	return stringBuilder.String()
}

func GenWorkflowYaml(workflowUsage, allNodeUsage string) string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(`根据要求生成一个工作流的配置文件
### 要求内容:
`)
	stringBuilder.WriteString(workflowUsage)
	stringBuilder.WriteString("\n\n")
	stringBuilder.WriteString(`### 各个节点的功能使用说明
`)
	stringBuilder.WriteString(allNodeUsage)
	stringBuilder.WriteString("\n")
	stringBuilder.WriteString(`
### 下面是一些工作流生成的节点参考
示例一:定时读取数据库数据，如果有新的数据开始执行读取数据调用chatgpt 节点生成描述信息，接下来调用tts节点将文本转为语音，并且把语音保存到S3上，接下来将语音文件地址保存到数据库中
name: ai-voice-location
storages:
  - name: qiniu
    type: s3
    bucket: xytschool  # 替换为你的S3存储桶名称
    region: cn-east-1         # 替换为你S3存储桶所在的区域
    access_key_id: xxx        # 替换为你的AWS访问密钥ID
    secret_access_key: xxx # 替换为你的AWS秘密访问密钥
    endpoint: https://s3.cn-east-1.qiniucs.com        # 可选，适用于使用AWS官方服务。如果使用的是自托管的S3兼容服务，可以替换为相应的endpoint
    path_style: false         # 可选，设置为true时使用路径风格访问，否则使用虚拟主机风格访问

nodes:
  - name: Webhook
    type: webhook
    parameters:
      httpMethod: POST
      path: webhook
      port: 8088
      exportHeader: # 支持将
        Authorization: "Authorization"
        Host: "Host"
      exportQuery: # 支持将
        rand: "rand"
      exportBody: # 支持将
        desc: "desc"
        locationName: "locationName"
        locationID: "locationID"
    global:
      locationName: "={{ $last.locationName }}"
      locationID: "={{ $last.locationID }}"
  - name: desc
    type: chatgpt
    credentials:
      openai_proxy: 'https://api.chatanywhere.tech/v1'
      openai_api_key: 'sk-xxx'
    parameters:
      question:  "=请帮忙生成一个景点的50字左右的描述信息(可以适当使用一些古诗词)：景区名称 {{ $last.locationName }}，景区简介{{ $last.desc }}"
      role: "assistant"
      roleDesc: "你是一名导游"
    global:
      generateContent: "={{ $last.content }}"
  - name: text2voice
    type: tts
    credentials:
      openai_proxy: 'https://api.chatanywhere.tech/v1'
      openai_api_key: 'sk-'
      baidu_api_key: ''
      baidu_secret: ''
    artifact:
      storage_name: qiniu
      from_path: "1"
      target_path: "=ai_gen/{{ $global.locationName }}_{{$global.locationID}}.mp3"
      is_overwrite: true
    parameters:
      text: "={{ $last.content }}"
      voice: "alloy"
    global:
      artifactKey: "=https://data.xytschool.com/ai_gen/{{ $global.locationName }}_{{$global.locationID}}.mp3"
  - name: save2Mysql
    type: mysql
    credentials:
      hostname: 'db.test.xxx'
      port: 'xx'
      database: 'xx'
      username: 'xx'
      password: 'xx'
    parameters:
      action: "update"
      sql: "=update locations set voice_url = ? ,info = ? WHERE id = ?"
      data:
        - '={{ $global.artifactKey }}'
        - '={{ $global.generateContent }}'
        - '={{ $global.locationID }}'
  - name: restartServer
    type: docker_compose
    parameters:
      action: "up"
      docker-compose-file: "={{ $last.newFile }}"
      workdir: "={{ $global.deployDir }}"
pinData: {}
global:
  gitUrl: "git.test.com"
  rootDir: /tmp
  deployDir: '/Users/xxx/code/go/src/gitlabee.chehejia.com/go-clean-template'
connections:
  Webhook:
      - - node: desc
  desc:
      - - node: text2voice
  text2voice:
      - - node: save2Mysql

示例二: webhook 触发器接收到请求后开始发布新的镜像，调用file节点修改docker-compose.yml中镜像的tag。调用docker_compose up 命令启动新的镜像，发送飞书消息到飞书群
name: deploy-image

nodes:
  - name: Webhook
    type: webhook
    credentials:
      frp:
        serverAddr: 'frp.lo.xxx'
        serverPort: 2023
        token: 'xxx!'
        customDomains:
          - 'deploylo.xxxx'
        localIP: "127.0.0.1"
        proxyName: 'workflow2'
    parameters:
      httpMethod: GET
      path: ui2
      port: 8101
      exportQuery: # 支持将
        image: image
        user: user
        deploy_env: env
    global:
      image: "={{ $last.image }}"
      deploy_env: "={{ $last.deploy_env }}"
  - name: changeReleaseVersion
    type: file
    parameters:
      action: "cp-replace"
      from: '(registry.cn-beijing.aliyuncs.com/xxxx):.+'
      to: "=${1}{{ $global.image }}"
      file: "={{ $global.deployDir }}/docker-compose.yml"
      workdir: "={{ $global.deployDir }}"
  - name: restartServer
    type: docker_compose
    parameters:
      action: "up"
      docker-compose-file: "={{ $last.distFile }}"
      workdir: "={{ $global.deployDir }}"
  - name: feishu_msg
    type: feishu_custom_robot
    credentials:
      webhook_url: 'https://open.feishu.cn/open-apis/bot/v2/hook/1ddc590c-0b5f-49a2-bb5f-8d00ade8ec4c'
    parameters:
      text:  "= {{ $last.content }}"
global:
  rootDir: /tmp
  deployDir: '/Users/mac/code/go/src/github.com/mytoolzone/workflow/demo/deploy-ui'

connections:
  Webhook:
    - - node: changeReleaseVersion
  changeReleaseVersion:
    - - node: restartServer
  restartServer:
    - - node: feishu_msg


`)
	return stringBuilder.String()
}
