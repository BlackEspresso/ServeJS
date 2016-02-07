function onRequest(resp,req){
	resp.contentType='text/html'
	//return JSON.stringify(req.cookies);
	router(resp,req)
		.on('/hello',hello)
		.on('/mailto', mailTo)
		.on('/run', run)
		.on('/login',login)
		.on('/admin',admin)
		.on('/req',request)
		.on('/editjs',editjs)
		.on('/boerse',boerse)

	addMapping('/websocket','websocket')
	addMapping('/writefile','writefile')
}

function boerse(resp,req){
    var symbol = req.formValues.symbol;
    var url = 'http://finance.yahoo.com/webservice/v1/symbols/'+symbol+'/quote?format=json';
    var siteResp = httpDo({url:url,method:'GET'});
    var serviceResp = JSON.parse(siteResp.body);
    //resp.write(JSON.stringify(JSON.parse(siteResp.body),null,2))
    resp.write(serviceResp.list.resources[0].resource.fields.name);
    resp.write('\n')
    resp.write(serviceResp.list.resources[0].resource.fields.price);
}

function editjs(resp,req){
	reloadTemplates();
	var js = readFile('js','main.js');
	var tmpl = runTemplate("EditJs.thtml",{MainJs:js.Suc});
	resp.write(tmpl)
}

function onWebSocketRequest(message){
	console.log('js: '+message)
}

function request(resp,req){
	siteResp = httpDo({url:'http://google.com',method:'GET'})
	resp.write(JSON.stringify(siteResp))
}

function mailTo(resp,req){
	var receiver = req.formValues.mailto
	var k = send(receiver,'hallo',JSON.stringify(req))
	resp.write(k.Err)
}

function run(resp,req){
	resp.contentType='text/plain'
	var cmd = runCmd('echo','4','5')
	if(cmd.Suc)
		resp.write(cmd.Suc);
	else
		resp.write(cmd.Err);
}

function hello(resp,req){
	resp.write('hello')
}

function login(resp,req){
	//reloadTemplates();
	loadMailSettings();
	if(req.method=='GET'){
		//resp.write()
		resp.write(JSON.stringify(settings))
		//resp.write("hello")
		//writeFile('static','test.html','<html>')
		//resp.write(resolve('m.yelp.com'))
		//resp.write(addTask('test',0,function(){console.log('ok')}))
		//startTasks()
		//resp.write(runTemplate("login.thtml",{Name:"v"}))
		//resp.write(JSON.stringify(req.formValues))
		//resp.write(JSON.stringify(readDir('static')))
	}
	else{
		//resp.write('logged in. no. kidding. password wrong')
	}
}

function admin(resp){
	resp.write("todo")
}

function router(resp,req){
	var obj= {
		on:function(url,func){
			if(req.url.indexOf(url)===0){func(resp,req)}
			return obj;
		}
	};
	return obj;
}






