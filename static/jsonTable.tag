<jObject>
  <h3>jObject</h3>
  <table>
    <tr each={k in kv}>
      <td>{k.name}</td>
      <td><raw content="{k}"/></td>
    </tr>
  </table>

  <script>
    this.data = opts.data || {}
    this.kv = [];

    this.on('mount', function() {
      for(var x in this.data){
        var val = this.data[x]
        tag = 'jValue'
        if(Array.isArray(val)){
          tag = 'jArray';
        }
        this.kv.push({name:x,val:val,tag:tag});
      }
      console.log(this.kv)
      this.update();
    })

  </script>
</jObject>

<jValue>
  <span>{val}</span>
  <script>
    this.val = opts.content.val||'';
    console.log('jvalue:' + this.val)
  </script>
</jValue>

<raw>
  <span></span>
 riot.mount(this.root, opts.content.tag, opts.content)
</raw>

<jArray>
</jArray>

<jTable>

</jTable>

<jRow>

</jRow>
