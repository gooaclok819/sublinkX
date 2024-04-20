<script setup lang='ts'>
import { ref,onMounted,nextTick  } from 'vue'
import {getNodes,AddNodes,DelNode,UpdateNode} from "@/api/subcription/node"
interface Node {
  ID: number;
  Name: string;
  Link: string;
  CreateDate: string;
}
const tableData = ref<Node[]>([])
const Nodelink = ref('')
const NodeOldlink = ref('')
const Nodename = ref('')
const NodeOldname = ref('')
const dialogVisible = ref(false)
const table = ref()
const NodeTitle = ref('')
const radio1 = ref('1')

async function getnodes() {
  const {data} = await getNodes();
    tableData.value = data
}
onMounted(async() => {
   getnodes()
})
const handleAddNode = ()=>{
  NodeTitle.value= '添加节点'
  Nodelink.value = ''
  Nodename.value = ''
  radio1.value = '1'
  dialogVisible.value = true

}
const addnodes = async ()=>{
  let nodelinks = Nodelink.value.split(/[\r\n,]/);
  // 分开 过滤空行和空格
  nodelinks = nodelinks.map((item) => item.trim()).filter((item) => item !== '');  
   if (NodeTitle.value== '添加节点'){
      // 判断合并还是分开
      if (radio1.value === '1') {
        if (Nodename.value.trim() === '') {
          ElMessage.error('备注不能为空')
          return
        }
        if (nodelinks) {
        Nodelink.value = nodelinks.join(',');
        await AddNodes({
        link: Nodelink.value.trim(),
        name: Nodename.value.trim(),
      })
      }
      } else {
        for (let i = 0; i < nodelinks.length; i++) {
          await AddNodes({
            link: nodelinks[i],
            name: "",
          })
        }
      }
      ElMessage.success("添加成功");
   }else{
    await UpdateNode({
        oldname: NodeOldname.value.trim(),
        oldlink: NodeOldlink.value.trim(),
        link: Nodelink.value.trim(),
        name: Nodename.value.trim(),
      })
    ElMessage.success("更新成功");
   }
    getnodes()
    Nodelink.value = ''
    Nodename.value = ''
    dialogVisible.value = false;
}

const multipleSelection = ref<Node[]>([])
const handleSelectionChange = (val: Node[]) => {
  multipleSelection.value = val
  
}
const selectAll = () => {
    nextTick(() => {
        tableData.value.forEach(row => {
            table.value.toggleRowSelection(row, true)
        })
    })
}
const handleEdit = (row:any) => {
  radio1.value = '1'
  for (let i = 0; i < tableData.value.length; i++) {
    if (tableData.value[i].ID === row.ID) {
      NodeTitle.value= '编辑节点'
      Nodename.value = tableData.value[i].Name
      NodeOldname.value = Nodename.value
      Nodelink.value = tableData.value[i].Link
      NodeOldlink.value = Nodelink.value
      dialogVisible.value = true
      // value1.value = tableData.value[i].Nodes.map((item) => item.Name)
    }
  }
}
const toggleSelection = () => {
  table.value.clearSelection()
}

const handleDel = (row:any) => {
  ElMessageBox.confirm(
    `你是否要删除 ${row.Name} ?`,
    '提示',
    {
      confirmButtonText: 'OK',
      cancelButtonText: 'Cancel',
      type: 'warning',
    }
  ).then(async () => {
      await DelNode({
        id: row.ID
      })
      ElMessage({
        type: 'success',
        message: '删除成功',
      })
      getnodes()
      // tableData.value = tableData.value.filter((item) => item.ID !== row.ID)
      
    })
  // console.log('click',row.ID)
}

const selectDel = () => {
  if (multipleSelection.value.length === 0) {
      return
  }
  ElMessageBox.confirm(
    `你是否要删除选中这些 ?`,
    '提示',
    {
      confirmButtonText: 'OK',
      cancelButtonText: 'Cancel',
      type: 'warning',
    }
  ).then( () => {
    for (let i = 0; i < multipleSelection.value.length; i++) {
       DelNode({
        id: multipleSelection.value[i].ID
      })
        tableData.value = tableData.value.filter((item) => item.ID !== multipleSelection.value[i].ID)
      }
      ElMessage({
        type: 'success',
        message: '删除成功',
      })
    })

}
// 分页显示
const currentPage = ref(1);
const pageSize = ref(10);
const handleSizeChange = (val: number) => {
  pageSize.value = val;
  // console.log(`每页 ${val} 条`);
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val;
  // console.log(`当前页: ${val}`);
}
// 表格数据静态化
const currentTableData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return tableData.value.slice(start, end);
});
// 复制链接
const copyUrl = (url: string) => {
  const textarea = document.createElement('textarea');
  textarea.value = url;
  document.body.appendChild(textarea);
  textarea.select();
  try {
    const successful = document.execCommand('copy');
    const msg = successful ? 'success' : 'warning';
    const message = successful ? '复制成功！' : '复制失败！';
    ElMessage({
      type: msg,
      message,
    });
  } catch (err) {
    ElMessage({
      type: 'warning',
      message: '复制失败！',
    });
  } finally {
    document.body.removeChild(textarea);
  }
};
const copyInfo = (row: any) => {
  copyUrl(row.Link)
}
</script>

<template>
  <div>
    <el-dialog
    v-model="dialogVisible"
    :title="NodeTitle"
    width="80%"
  >
  <el-input 
  v-model="Nodelink" 
  placeholder="请输入节点多行使用回车或逗号分开,支持base64格式的url订阅" 
  type="textarea" 
  style="margin-bottom:10px" 
  :autosize="{ minRows: 2, maxRows: 10}"
  />
  <el-radio-group v-model="radio1" class="ml-4" v-if="NodeTitle== '添加节点'">
      <el-radio value="1" size="large">合并</el-radio>
      <el-radio value="2" size="large">分开</el-radio>
    </el-radio-group>
  <el-input v-model="Nodename" placeholder="请输入备注"  v-if="radio1!='2'" />
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="addnodes">确定</el-button>
      </div>
    </template>
  </el-dialog>
    <el-card>
    <el-button type="primary" @click="handleAddNode">添加节点</el-button>
    <div style="margin-bottom: 10px"></div>
      <el-table ref="table" :data="currentTableData" style="width: 100%" @selection-change="handleSelectionChange">
    <el-table-column type="selection" fixed prop="ID" label="id"  />
    <el-table-column prop="Name" label="备注"  >
      <template #default="scope">
        <el-tag type="success">{{scope.row.Name}}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="Link" label="节点" sortable :show-overflow-tooltip="true" />
    <el-table-column prop="CreateDate" label="创建时间" sortable  />
    <el-table-column fixed="right" label="操作" width="120">
      <template #default="scope">
        <el-button link type="primary" size="small" @click="copyInfo(scope.row)">复制</el-button>
        <el-button link type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>
  <el-button link type="primary" size="small" @click="handleDel(scope.row)">删除</el-button>

      </template>
    </el-table-column>
  </el-table>
  <div style="margin-top: 20px" />
    <el-button type="info" @click="selectAll()">全选</el-button>
    <el-button type="warning" @click="toggleSelection()">取消选择</el-button>
    <el-button type="danger" @click="selectDel">批量删除</el-button>
  <div style="margin-top: 20px"/>
  <el-pagination
  @size-change="handleSizeChange"
  @current-change="handleCurrentChange"
  :current-page="currentPage"
  :page-size="pageSize"
  layout="total, sizes, prev, pager, next, jumper"
  :page-sizes="[10, 20, 30, 40]"
  :total="tableData.length">
</el-pagination>

    </el-card>
  </div>
</template>

<style scoped>
.el-card{
  margin: 10px;
}
.el-input{
  margin-bottom: 10px;
}
</style>