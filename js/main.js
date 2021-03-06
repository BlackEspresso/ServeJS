function onStart(){
	var m = require('httpmappings')
	var t = require('tasks')
	var c = require('cache')
	
	m.addMapping('/websocket','websocket');
	m.addMapping('/writefile','writefile');
	m.addMapping('/static','servefile');
	t.startTasks();
	
	var yaml = require('settings');
	var conf = yaml.read('./serverjs.yaml');
	c.set('settings', JSON.stringify(conf.ok))

	var server = require('httplistener')
	server.requestFuncName = 'onRequest';
	//server.addr = ':8080'
	server.startAndWait()
}

function onRequest(resp,req){
	resp.contentType='text/html'
	//return JSON.stringify(req.cookies);
	router(resp,req)
		.on('/gquery',gquery)
		.on('/echo',echo)
		.on('/hello',hello)
		.on('/showcache',showcache)
		.on('/mailto', mailTo)
		.on('/mm',runa)
		.on('/run', runCmd)
		.on('/scanandmail',scanandmail)
		.on('/login',login)
		.on('/register',register)
		.on('/req',request)
		.on('/editjs',editjs)
		.on('/boerse',boerse)
		.on('/scandns',scandns)
		.on('/schedule',schedule)
		.on('/datetime',datetime)
		.on('/cache',cacheFunc)
		.on('/header',header)
		.on('/htmlcheck',htmlCheck)
		.on('/mongo',mongo)
		.on('/siteinfo',siteInfo)
		.on('/sitescan',siteScan)
		.on('/time',time)
		.on('/userfiles',userfiles)
		.on('/testHttp',testHttp)
		.on('/filewatch',filewatch)
		.on('/sendwebsocket',sendwebsocket)
		.on('/stats',stats)
		.on('/xsstest',xsstest)
		
}


function getQuery(url){
	var splitted = url.split('?');
	var queries = {};
	if(splitted.length<=1){
		return queries;
	}

	var queryString = splitted[1];
	var qs = queryString.split('&');
	for(var q=0;q<qs.length;q++){
		var kv = qs[q].split('=');
		queries[kv[0]]=kv[1];
	}
	return queries;
}

function scanandmail(resp,req){
	var tasks= require('tasks');
	var url = req.formValues.url[0];
	var email = req.formValues.email[0];
	
	tasks.addTask('test',0,function(){
		var baseUrl = url;
		var vecs = ['><zq>','\"><zq>','\'><zq>'];
		var queries = getQuery(baseUrl);

		var mt = '';
		var log = {};
		var founds = [];

		var baseReq = getSiteInfo(baseUrl);
		founds.push(baseReq);

		for(var q in queries){
			for(var x=0;x<vecs.length;x++){
				var nextUrl = baseUrl.replace(q+'='+queries[q],q+'='+vecs[x]+queries[q])
				console.log(nextUrl)
				var ret = getSiteInfo(nextUrl);
				if(ret.zsq.length >0){
					founds.push(nextUrl)
				}
				log.push(ret);
			}
		}

		mt += JSON.stringify(founds);
		mt += '\n'
		for(var x=0;x<log.length;x++){
			mt += log[x].request.url + ' ' + log[x].invalidTags.length;
			mt += '\n'
		}
		


		var mail = require('mail');
		var err = mail.send(email,'scanandmail', mt);
	});
}

function xsstest(resp,req){
	resp.write('<html>')
	var link = req.formValues.link|| 'test';
	resp.write('<a href="'+link+'">'+ link+'</a>')
	resp.write('</html>')
}

function stats(resp,req){
	var stats = require('stats');
	resp.write(JSON.stringify(stats.getStats()));
}

function runa(resp,req){
	var modules = require('modules')
	var time = require('time')
	var k = 45;
	modules.run(function(){console.log(k)})
	time.sleep(10000)
	resp.write('ok')
}

function sendwebsocket(resp,req){
	var events = require('events');
	events.push('channel1','hello from http')
}

function filewatch(resp,req){
	var fw = require('filewatch')
	fw.watchDir('./static/')
	fw.start()
}

function testHttp(resp,req){
	var http = require('http');
	var url = 'http://localhost:8081/echo'
	var cResp = http.do({
		method:'POST',
		url:url,
		header:{
			'Test':'4',
			'User-Agent':'QQ',
			'Cookie':'test=4',
			'Content-Type':'application/x-www-form-urlencoded'
		},
		body:'id=5&name=hello'
	});

	if(cResp.error)
		resp.write(cResp.error)
	else
		resp.write(cResp.ok.body)

}

function userfiles(resp,req){
	var id = req.formValues.id;
	if(id==null){
		resp.statusCode = 404;
		return;
	}
	var r = /\d+/;
	if(!r.test(id)){
		resp.statusCode = 404;
		return;
	}
	var file = require('file');
	var c = file.read('./userfiles/reports/'+id+'.txt');
	if(c.error!==undefined){
		resp.statusCode = 404;
		return;
	}

	resp.write(c.ok);
}

function time(resp,req){
	var time = require('time');
	time.sleep(4000)
	resp.write(4000);
}

function siteScan(resp,req){
	var urls = req.formValues.url;
	if(urls==null||urls.length==0){
		resp.write('url not found')
	}
	var url = urls[0];
	if(url[url.length-1]=="/"){
		url = url.substr(0,url.length-1);
	}
	var count = req.formValues.count;
	if(count == null)
		count = 10;
	var start = req.formValues.start;
	if(start == null)
		start = 0;
	var email = req.formValues.email;
	if(email == null){
		resp.write('email is empty');
		return;
	}

	var tasks= require('tasks');

	tasks.addTask('test',0,function(){
        var mail= require('mail');
        var ret = fuzzUrls(url,count,start);

        var fileId = (Math.random()*100000)|1;
        var file = require('file');
        var fileName = './userfiles/reports/'+fileId+'.txt';
        file.write(fileName,ret);
        var fileUrl = 'http://' + req.host + '/userfiles?id='+fileId;
        var nextRequestUrl = 'http://' + req.host + req.url;
        var err = mail.send(email,'scanurl',fileUrl + '\n\n' + nextRequestUrl);
        console.log(err.error);

    });

	resp.write('ok');
}

function fuzzUrls(url,count,start){
	var http = require('http');
	var file = require('file');
	var wordlist = file.read('./static/wordlists/KitchensinkDirectories.fuzz.txt');
	var words = wordlist.ok.split('\n');

	var ret = '';

	var startDate = new Date();
	var time = require('time');

	for(var x = start;x<words.length&&x-start<count;x++){
		var fullUrl = url + words[x].trim();
		var cResp = http.do({url:fullUrl});
		time.sleep(1000);
		ret += cResp.ok.statusCode + '  ' + fullUrl + '\n';
	}
	var durationSec = (new Date().getTime()- startDate.getTime())/1000;
	ret += durationSec + ' seconds\n'
	return ret;
}

function getSiteInfo(url){
	 var goquery = require('goquery')
	var http = require('http')
	var cResp = http.do({
		url:url,
		followRedirects:false,
		header:{
			'Accept':'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
			'Accept-Language':'de,en-US;q=0.8,en;q=0.6',
			'User-Agent':'Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36'
		}
	});
	
	if (cResp.error!=null){
		console.log(cResp.error)
		resp.write(cResp.error)
		return
	}
	//console.log(cResp.ok.body)
	var doc = goquery.newDocument(cResp.ok.body)
	var form = doc.ExtractForms();
	var hrefs = doc.ExtractHrefs(url);
	var scripts = doc.ExtractAttributes('script');
	var links = doc.ExtractAttributes('link');
	var zsq = doc.ExtractAttributes('zq');
	
	var htmlcheck = require('htmlcheck')
	var k = htmlcheck.loadTags('./static/tags.json')
	//console.log(k.error)
	var err = htmlcheck.validate(cResp.ok.body)

	var reg = /((https?|ftp|file):)?\/\/[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[a-zA-Z0-9+&@#/%=~_|]/g;
	var textUrls = cResp.ok.body.match(reg);

	if (textUrls==null){
		textUrls = [];
	}

	var hrefsObj = [];
	for(var x=0;x<hrefs.length;x++){
		hrefsObj.push({href:hrefs[x]})
		textUrls.push(hrefs[x]);
	}
	
	textUrls.sort();
	var allUrls = [];
	var c = {}
	for(var x=0;x<textUrls.length;x++){
		if(c[textUrls[x]]!=true){
			c[textUrls[x]]=true;
			allUrls.push({href:textUrls[x]});
		}
	}


	var ret = {
		request:{ url:url},
	    zsq:zsq,
	    invalidTags:err,
	    header:cResp.ok.header,
	    hrefs:hrefsObj,
	    forms:form,
	    allUrls:allUrls,
	    scripts:scripts,
	    links:links,
	    cookies:cResp.ok.cookies,
	    status:cResp.ok.statusCode
	    tls: cResp.ok.tlsDNSNames
	    
	}
	return ret;
}

function siteInfo(resp,req){
    var url = req.formValues.url[0];
   	var ret = getSiteInfo(url);

	resp.write(JSON.stringify(ret))
}

function showcache(resp,req){
	var cache = require('cache');
	if(!isAuth(req)){
		redirect(resp,'/login')
		return
	}
	resp.write(JSON.stringify(cache.all()))
}

function redirect(resp,url){
	resp.statusCode=302;
	resp.header = {
		'Location':url
	}
}

function mongo(resp,req){
	var m = require('mongodb')
	var sess = m.newSession('localhost')
	var c= sess.DB('jstest').C('test')

	c.RemoveAll({})

	c.Insert({'Name':'Tester1','Street':'here'})

	c = c.Find({'Name':'Tester1'})
	for(var x in c){
		resp.write(x + '\n');
	}

	var test = m.all(c.Select({'Name':1}));

	resp.write(JSON.stringify(test));
}

function gquery(resp,req){
	var goquery = require('goquery')
	var http = require('http')
	var cResp = http.do({url:'http://google.com'});
	var doc = goquery.newDocument(cResp.ok.body)
	var links = doc.ExtractAttributes('body');
	resp.write(JSON.stringify(links))
}

function htmlCheck(resp,req){
    var url = req.formValues.url;
    if(url==null){
        resp.write('url is empty')
        return
    }
    var cResp = http.do({url:url});
    htmlcheck.loadTags('./static/tags.json')
    
	var err = htmlcheck.validate(cResp.ok.body)
	resp.write(JSON.stringify(err))
}

function header(resp,req){
	resp.header = {
		'Set-Cookie':'c=5',
		'Test-Header':'4'
	}
	resp.write('4')
}

function cacheFunc(resp,req){
	var key = req.formValues.key;
	var val = req.formValues.okue;
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
    var serviceResp = JSON.parse(siteResp.ok.body);
    //resp.write(JSON.stringify(JSON.parse(siteResp.body),null,2))
    resp.write(serviceResp.list.resources[0].resource.fields.name);
    resp.write('\n')
    resp.write(serviceResp.list.resources[0].resource.fields.price);
}

function scandns(resp,req){
	var file = require('file')
    var fileRead = file.read('static', 'subdomain_wordlist.txt');
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
	var templ = require('templating');
	if(!isAuth(req)){
		redirect(resp,'/login');
		return;
	}
	var file = require('file')
	templ.reloadTemplates();
	var js = file.read('./js/main.js');
	var tmpl = templ.runTemplate("EditJs.thtml",{MainJs:js.ok});
	resp.write(tmpl)
}

function onWebSocket(websocket){
	var e = require('events');
	var myChannel = 'webrandom1';
	e.create(myChannel)
	e.route('channel1',myChannel)
	while(true){
		var message = e.next(myChannel);
		websocket.write(1,message)
		e.sleep(1000)
		break
	}
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

function runCmd(resp,req){
	resp.contentType='text/plain'
	var cmd = runCmd('echo','4','5')
	if(cmd.ok)
		resp.write(cmd.ok);
	else
		resp.write(cmd.error);
}

function echo(resp,req){
	resp.write(JSON.stringify(req));
}

function hello(resp,req){
	var c = require('crypto');
	resp.write('hello ' + c.newGuid());
	//resp.write('hello')
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

function getCollection(mongodb){
	var sess = mongodb.newSession('localhost')
	var c = sess.DB('jstest').C('test')
	return c;
}

function register(resp,req){
	var mongodb = require('mongodb')
	var c = getCollection(mongodb);

	var templating = require('templating');
	templating.reloadTemplates();

	if(req.method=='GET'){
		var templ = templating.runTemplate("login.thtml",{ActionName:'register'});
		resp.write(templ)
	}
	if(req.method=='POST'){
		var username = req.formValues.username[0];
		var password = req.formValues.password[0];
		c.Insert({'Username':username,'Password':password});
	}
}

function getCookie(req,name){
	var cookies = {};
	if(req.cookies!=null){
		for(var x = 0;x<req.cookies.length;x++){
			if(req.cookies[x].Name==name)
				return req.cookies[x];
		}
	}
	return undefined;
}

function getAuth(req){
	var cache = require('cache');
	var cookie = getCookie(req,'userid')
	if (cookie == null)
		return null;
	var userId = cookie.Value;
	var fromCache = cache.get('userid_'+userId);
	if(fromCache == null)
		return null;
	var obj = JSON.parse(fromCache);
	return obj;
}

function isAuth(req){
	var auth = getAuth(req);
	return auth!==null && auth.isAuth;
}

function login(resp,req){
	var templating = require('templating');
	templating.reloadTemplates();

	if(req.method=='GET'){
		var auth = getAuth(req);
		if(auth != null && auth.isAuth){
			message = 'hello ' + auth.username
		} else{
			message = 'please log in'
		}

		var templ = templating.runTemplate("login.thtml",
			{ActionName:'login',Message:message});
		resp.write(templ)
	}
	else{
		var mongodb = require('mongodb')
		var cache = require('cache')
		var c = getCollection(mongodb);

		var username = req.formValues.username[0];
		var password = req.formValues.password[0];

		c =c.Find({'Username':username})
		var user = mongodb.one(c);
		if(user !== false && user.Password==password){
			var auth = getAuth(req);
			if(auth!=null && auth.id!=null){
				cache.remove('userid_'+auth.id)
			}
			var id = (Math.random()*100000000)|1;

			cache.set('userid_'+id, JSON.stringify({
				isAuth:true,
				id:id,
				date: new Date(),
				isAdmin:user.isAdmin === true,
				username:username})
			);
			resp.header = {'Set-Cookie':'userid='+id+';HttpOnly'};

		}
		else
			resp.write('not ok');
	}
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






