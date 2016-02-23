<jobject>
  <table>
    <tbody>
    <tr each={k in kv}>
      <td>{k.name}</td>
      <td><raw content="{k}"/></td>
    </tr>
  </tbody>
  </table>

  <script>
    this.data = opts.data || {}
    this.kv = [];

    this.on('update', function() {
      for(var x in this.data){
        var val = this.data[x]
        var tag = 'jvalue';
        if(Array.isArray(val)){
          tag = 'jtable';
        }else if(typeof val=="object"){
          tag ='jobject';
        }
        this.kv.push({name:x,data:val,tag:tag});
      }
      this.update();
    })

  </script>
</jobject>

<jtable>
  <table>
    <thead>
    <tr>
      <th each="{row in columns}">{row}</th>
    </tr>
  </thead>
  <tbody>
    <tr each="{row in rows}" if={!isArray}>
      <td each="{k,v in row}">{v}</td>
    </tr>
    <tr if={isArray}>
      <td each="{k,v in rows}">{k}</td>
    </tr>
  </tbody>
  </table>
  <script>
    this.rows = opts.data || {}
    this.columns = [];
    this.isArray=false;
    this.on('update', function() {
      var firstElement = this.rows[0];
      if (typeof firstElement != 'object'){
        this.isArray=true;
      } else {
        for(var x in firstElement)
            this.columns.push(x)
      }
      this.update();
    });
  </script>
</jtable>

<jvalue>
  <span>{val}</span>
  <script>
    this.val = opts.data||'';
  </script>
</jvalue>

<raw>
 riot.mount(this.root, opts.content.tag, opts.content)
</raw>
