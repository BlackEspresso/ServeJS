function onStart(){
	var m = require('httpmappings')
	var t = require('tasks')
	var c = require('cache')
	
	m.addMapping('/websocket','websocket');
	m.addMapping('/writefile','writefile');
	t.startTasks();
	
	var yaml = require('settings');
	var conf = yaml.read('./serverjs.yaml');
	c.set('settings', JSON.stringify(conf))

	var server = require('httplistener')
	server.requestFuncName = 'onRequest';
	//server.addr = ':8080'
	server.startAndWait()
}

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
		.on('/scandns',scandns)
		.on('/schedule',schedule)
		.on('/datetime',datetime)
		.on('/cache',cacheFunc)
		.on('/header',header)
		.on('/htmlcheck',htmlCheck)
}

function htmlCheck(resp,req){
    var url = req.formValues.url;
    if(url==null){
        resp.write('url is empty')
        return
    }
    var cResp = http.do({url:url});
    htmlcheck.loadTags('./static/tags.json')
    
	var err = htmlcheck.validate(cResp.val.body)
	resp.write(JSON.stringify(err))
}

function header(resp,req){
	resp.header = {
		'Set-Cookie':'c=5',
		'Test-Header':'4'
	}
}

function cacheFunc(resp,req){
	var key = req.formValues.key;
	var val = req.formValues.value;
	if (val == null){
		resp.write(cache.get(key))
	}else{
		cache.set(key,val);
		resp.write('ok')
	}
	cache.save('./static/cache.json')
} 

function datetime(resp,req){
	resp.write(new Date().toString())
}

function boerse(resp,req){
    var symbol = req.formValues.symbol;
    var url = 'http://finance.yahoo.com/webservice/v1/symbols/'+symbol+'/quote?format=json';
    var siteResp = http.do({url:url});
    //resp.write(JSON.stringify(siteResp))
    var serviceResp = JSON.parse(siteResp.val.body);
    //resp.write(JSON.stringify(JSON.parse(siteResp.body),null,2))
    resp.write(serviceResp.list.resources[0].resource.fields.name);
    resp.write('\n')
    resp.write(serviceResp.list.resources[0].resource.fields.price);
}

function scandns(resp,req){
    var fileRead = readFile('static', 'subdomain_wordlist.txt');
    var lines = fileRead.Suc.split('\n');

    if (req.formValues.host === null) {
        resp.write('error')
    }

    var count = 5;
    if (req.formValues.count !== null) {
        count = parseInt(req.formValues.count)
    }

    var ret = {};
    var regExIp = /([^\s]+)\s+(\d+)\s+(\w+)\s+(\w+)\s+(\d+\.\d+\.\d+\.\d+)/ig;

    resp.write('[\n')
    for (var x = 0; x < lines.length&&x<count; x++) {
        var dnsReponse = resolve(lines[x].trim() + '.' + req.formValues.host);
        var entries= dnsReponse.split('\n');
        
        for(var z=0;z<entries.length;z++){
    		var matches= entries[z].split(/\s+/);
        	for (var y = 0; y < matches.length; y++) {
            	if (y > 4)
                	matches[4] += matches[y];
            }
        	matches.splice(5, matches.length - 5);
        	if(matches.length === 1 && matches[0]==="")
        		continue;
            resp.write(JSON.stringify(matches));
            if(x+1<lines.length && x+1<count){
                resp.write(',\n');
            }
        }
    }
    resp.write('\n]');
}

function editjs(resp,req){
	reloadTemplates();
	var js = readFile('./js/main.js');
	var tmpl = runTemplate("EditJs.thtml",{MainJs:js.val});
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
	if(cmd.val)
		resp.write(cmd.val);
	else
		resp.write(cmd.error);
}

function hello(resp,req){
	resp.write('hello')
}

function schedule(resp,req){
	var d = new Date();
	d.setMinutes(d.getMinutes()+1);

    addTask('test',d,function(){
        //var cmd = runCmd('cmd','/c','dir');
        console.log('task started: ' + new Date())
    });
	
	resp.write('scheduled')
}

function login(resp,req){
	loadMailSettings();
	if(req.method=='GET'){
		resp.write(runTemplate("login.thtml",{Name:"v"}))
		//resp.write(JSON.stringify(readDir('static')))
		//writeFile('static','test.html','<html>')
	}
	else{
		resp.write('logged in. no. kidding. password wrong')
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






