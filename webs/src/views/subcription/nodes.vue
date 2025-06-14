<script setup lang='ts'>
import { ref,onMounted,nextTick  } from 'vue'
import {getNodes,AddNodes,DelNode,UpdateNode,GetGroup,SetGroup} from "@/api/subcription/node"
import type { ElTable } from 'element-plus'

interface GroupNode {
  ID: number;
  Name: string;
  Nodes :Node[];
}
interface Node {
  ID: number;
  Name: string;
  Link: string;
  CreateDate: string;
  GroupNodes?: GroupNode[]; // 分组信息
  
}
interface NodeInfo {
    ID?: number // 编辑时需要传入ID
    Title?:string 
    Name?: string
    Link: string
    GroupName?: string[] // 分组名称
}
onMounted(async() => {  // 页面开始执行函数
   getnodes()
   GetGroups()
})
const dialogMode = ref<'add' | 'edit'>('add');

// --- 表格选择与操作相关数据 ---
const multipleSelection = ref<Node[]>([]); // Stores selected table items
const multipleTable = ref<InstanceType<typeof ElTable> | null>(null)


const tableRefs = ref<{ [key: string]: any }>({}); // Stores references to each el-table
// --- 表格选择与操作相关数据结束 ---
// const NodeNewLinkInput = ref("")
// const NodeNewNameInput = ref("")
const NodeGroupInput = ref("")
const tableData = ref<Node[]>([])
// 分组列表临时存放数据
const tableDataTemp = ref<Node[]>([])
// 分组列表临时存放数据
const activeName = ref('全部')
const Nodedialog = ref (false); // 弹窗是否可见
const Groupdialog = ref (false); // 弹窗是否可见
const NodeForm = ref<NodeInfo>({
    Title: '',
    Name: '',
    Link: '',
    GroupName: [],
  })
const allGroupNames = ref<string[]>([]); // 所有分组名称
const allNodes = ref<string[]>([]); // 所有节点
const nodelistShow = ref(false); // 节点列表
const SelectionNodeGroups = ref<string[]>([]); // 选中的分组
const SelectionNode = ref(''); // 选中的节点

// const SelectionNodes = ref([]); // 选中的节点
const RadioGroup = ref("1"); // 分组单选框
// 将所有输入的值清空
function ClearInput() {
  SelectionNode.value = ''; // 清空选中的节点
  NodeForm.value = { // 清空节点链接输入框
    Title: '',
    Name: '',
    Link: '',
    GroupName: [],
  }
  NodeGroupInput.value = ''; // 清空创建分组输入框
  SelectionNodeGroups.value = []; // 清空选中的分组
  nodelistShow.value = false; // 隐藏节点列表
  Nodedialog.value = false; // 关闭节点添加弹窗
  Groupdialog.value = false; // 关闭分组绑定弹窗
  
}
async function getnodes() {
  const {data} = await getNodes();
  if (data.length > 0) tableDataTemp.value = tableData.value = data
  allNodes.value = []; // 清空 allNodes 数组
  data.forEach((item:any) => {
      allNodes.value.push(item.Name); // 将所有节点添加到 allNodes 中
  });
  
} 
async function GetGroups() {
  const {data} = await GetGroup();
  if (Array.isArray(data) && data.length > 0) {
  allGroupNames.value=data; // 将所有分组名称添加到 allGroupNames 中

}
  RadioGroup.value = allGroupNames.value.length > 0 ? "1" : "2"; // 自动选择单选框值
  // console.log("单选框",RadioGroup.value);
  
}


const handleAddNode = () => {
  dialogMode.value = 'add';
  Nodedialog.value = true;
  NodeForm.value = {
    Title: '添加节点',
    Name: '',
    Link: '',
    GroupName: [],
  };
  SelectionNodeGroups.value = [];
  NodeGroupInput.value = '';
};

const handleEditNode = (row: Node) => {  
  // NodeNewNameInput.value = row.Name; // 编辑时使用原名称
  // NodeNewLinkInput.value = row.Link; // 编辑时使用原链接
  dialogMode.value = 'edit';
  Nodedialog.value = true;
  NodeForm.value = {
    ID: row.ID,
    Title: '编辑节点',
    Name: row.Name,
    Link: row.Link,
    GroupName: (row.GroupNodes || []).map(g => g.Name),
  };
  SelectionNodeGroups.value = NodeForm.value.GroupName || [];
  SelectionNode.value = row.Name;
};
const SubmitNodeForm = async (row:any) => {
  const isAdd = dialogMode.value === 'add';
  let links = NodeForm.value.Link.trim().split(/[\n,]/).map(item => item.trim()).filter(item => item);
  if (isAdd && links.length === 0) {
    ElMessage.warning('节点链接不能为空');
    return;
  }

  try {
    if (isAdd) {
      for (const link of links) {
        await AddNodes({
          link,
          group: RadioGroup.value === '1' ? SelectionNodeGroups.value.join(',') : NodeGroupInput.value,
        });
      }
      ElMessage.success('节点添加成功');
    } else {
      await UpdateNode({
        id:NodeForm.value.ID,
        name: NodeForm.value.Name, // 新名称
        link: NodeForm.value.Link, // 新链接
        group: RadioGroup.value === '1' ? SelectionNodeGroups.value.join(',') : NodeGroupInput.value,
      });
      ElMessage.success('节点更新成功');
    }


  } catch (err) {
    ElMessage.error(`${isAdd ? '添加' : '更新'}失败`);
  }
  getnodes();
  GetGroups();
  ClearInput();
};

// const AddNode = async() => {
//   // 多节点链接输入处理
//   let NodeLinkInputs = NodeNewLinkInput.value.trim().split(/[\n,]/); // 使用换行符或逗号分隔输入的节点链接
//   NodeLinkInputs = NodeLinkInputs.map((item) => item.trim()).filter((item) => item !== ''); // 去除空白和重复的链接
//   if (NodeNewLinkInput.value.trim() === '') {
//     ElMessage.warning('节点链接不能为空');
//     return;
//   }

//   try {
//     // 多节点同步循环添加节点
//     for(const link of NodeLinkInputs) {
//       if (link) {
//           const newNode = {
//           link: link.trim(), // 节点链接
//           group: SelectionNodeGroups.value.join(','), // 选中的分组
//           };
//           await AddNodes(newNode).then(() => {
//           ElMessage.success('节点添加成功');
//           Nodedialog.value = false; // 关闭弹窗
//           });
//       }
//     }
//     // getnodes(); // 刷新节点列表
//     // GetGroups(); // 刷新分组列表
//   } catch (error) {
//     console.error('添加节点失败:', error);
//     // ElMessage.error('添加节点失败，请稍后再试');
//   }
//   getnodes(); // 刷新节点列表
//   GetGroups(); // 刷新分组列表
//   ClearInput(); // 清空所有输入
// }
const AddGroup = async() => {
  console.log(SelectionNode.value);

  try {
    // 检查是否选择了已有分组或输入了新分组名
    console.log(RadioGroup.value, SelectionNodeGroups.value, NodeGroupInput.value);
    
    if (RadioGroup.value === "1" && SelectionNodeGroups.value.length === 0) {
      ElMessage.warning('你还没有选择分组');
      return;
    }
    if (RadioGroup.value === "2"&&NodeGroupInput.value.trim() === '') {
      ElMessage.warning('创建的分组名不能为空');
      return;
  }
      if (SelectionNode.value.length > 0) { // 如果没有选择节点
      const newNode = {
      name: SelectionNode.value, // 节点链接
      group: RadioGroup.value == '1' ?SelectionNodeGroups.value.join(','):NodeGroupInput.value, // 条件选择已有节点|创建分组
      };
      await SetGroup(newNode).then(() => {
      ElMessage.success('分组绑定成功');
            });
    }
  } catch (error) {
    console.error('添加分组失败:', error);
    // ElMessage.error('添加分组失败');
  }
  getnodes(); // 刷新节点列表
  GetGroups(); // 刷新分组列表
  ClearInput(); // 清空所有输入
}
// 表格时间格式化
const Timeformatter  = (row:any)=>{
  row.CreatedAt = new Date(row.CreatedAt).toLocaleString(); // 转换为本地时间字符串
  return row.CreatedAt;
  
}
// 选择已有节点显示所属分组
const  handleShownodeGroupList =()=>{
  // 显示这个节点关联的分组
  const nodeData = allNodes.value.find(node => node === SelectionNode.value);
  SelectionNodeGroups.value = []
  tableData.value.forEach((item, ) => {
    if (item.Name === nodeData && (item.GroupNodes?.length ?? 0) > 0) {
      // console.log(`节点 ${nodeData} 的分组:`, item.GroupNodes);
      item.GroupNodes?.forEach((item) => {
        SelectionNodeGroups.value.push(item.Name); // 将分组名称添加到 SelectionNodeGroups 中
      });
    } 
});
}
// 表格所属分组格式化
const Groupformatter = (row:any,cellValue:any) =>{
  const data = row.GroupNodes || [];
  if (!Array.isArray(data) || data.length === 0) {
    return '未分组'; // 如果没有分组，返回默认值
  }
 return data.map((group: any) => group.Name).join(', ');
}
// --- 复制链接 (保持不变) ---
const copyUrl = (url: string) => {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(url)
      .then(() => {
        ElMessage.success('链接已复制到剪贴板！');
      })
      .catch(err => {
        console.error('复制失败:', err);
        ElMessage.error('复制失败！请手动复制。');
      });
  } else {
    const textarea = document.createElement('textarea');
    textarea.value = url;
    document.body.appendChild(textarea);
    textarea.select();
    try {
      document.execCommand('copy');
      ElMessage.success('链接已复制到剪贴板！');
    } catch (err) {
      ElMessage.warning('复制失败！');
    } finally {
      document.body.removeChild(textarea);
    }
  }
};
// 复制表格节点信息
const copyInfo = (row: Node) => {
  copyUrl(row.Link);
};
const handleDel = async (row: Node) => {
  try {
    await ElMessageBox.confirm(
      `你是否要删除 ${row.Name} ?`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    );
    await DelNode({ id: row.ID });
    ElMessage.success('删除成功');
  } catch (error) {
    if (error !== 'cancel') {
      console.error("删除失败:", error);
      ElMessage.error('删除失败！');
    }
  }
  // 刷新节点列表
  await GetGroups(); // 刷新分组列表
  await getnodes(); // 刷新节点列表
  ClearInput(); // 清空所有输入
};
const selectDel = async () => {
  
  if (multipleSelection.value.length === 0) {
    ElMessage.warning('请选择要删除的节点！');
    return;
  }
  try {
    await ElMessageBox.confirm(
      `你是否要删除选中的 ${multipleSelection.value.length} 条节点 ?`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    );

    const IDs: number[] = []

    for (const item of multipleSelection.value) {
      await DelNode({ id: item.ID });
       IDs.push(item.ID); // 收集所有已删除的节点ID
    }
    ElMessage.success('批量删除成功');
    // 从 tableData 中删除已删除的节点
    tableData.value = tableData.value.filter(item => !IDs.includes(item.ID));
; 

  } catch (error) {
    if (error !== 'cancel') {
      console.error("批量删除失败:", error);
      ElMessage.error('批量删除失败！');
    }
  }
    // 刷新节点列表
  await GetGroups(); // 刷新分组列表
  await getnodes();
};
// 全选
const selectAll = () => {
  nextTick(() => {
const table = multipleTable.value
  if (table) {
      // 否则全选
      tableData.value.forEach(row => {
        table.toggleRowSelection(row, true)
      })
  }
  });
};
// 取消全选 
const selectClear = () => {
  nextTick(() => {
    const table = multipleTable.value;
    if (table) {
      table.clearSelection();
    }
  });
};
// --- 表格选择操作 (保持不变) ---
const setTableRef = (el: any, name: string) => {
  if (el) {
    tableRefs.value[name] = el;
  } else {
    delete tableRefs.value[name];
  }
};
//批量复制
const selectCopy = async () => {
  if (multipleSelection.value.length === 0) {
    ElMessage.warning('请选择要复制的节点！');
    return;
  }
  try {
    copyUrl(multipleSelection.value.map(item => item.Link).join('\n'));
  } catch (error) {
    if (error !== 'cancel') {
      console.error("批量复制失败:", error);
      ElMessage.error('批量复制失败');
    }
  }
};
const handleSelectionChange = (val: Node[]) => {
  multipleSelection.value = val;
};

watch(activeName, (newVal) => {
  if (newVal === '全部') {
    tableData.value = tableDataTemp.value;
  } else {
    tableData.value = tableDataTemp.value.filter(item => {
      return item.GroupNodes?.some(group => group.Name === newVal);
    });
  }
});


</script>

<template>
  <div>
 <el-dialog v-model="Nodedialog" :title="NodeForm.Title" width="80%">
  <el-input
    v-model="NodeForm.Link"
    placeholder="请输入节点链接，支持多行使用回车或逗号分开"
    type="textarea"
    style="margin-bottom: 10px"
    :autosize="{ minRows: 2, maxRows: 10 }"
    v-if="dialogMode === 'add'"
  />

<el-input
  v-model="NodeForm.Name"
  placeholder="节点名称（编辑时）"
  style="margin-bottom: 10px"
  v-if="dialogMode === 'edit'"
/>
  <el-input
    v-model="NodeForm.Link"
    placeholder="请输入节点链接，支持多行使用回车或逗号分开"
    type="textarea"
    style="margin-bottom: 10px"
    :autosize="{ minRows: 2, maxRows: 10 }"
    v-if="dialogMode === 'edit'"
  />

  <!-- 分组部分 -->
  <el-radio v-model="RadioGroup" label="1" v-if="allGroupNames.length > 0">选择已有分组</el-radio>
  <el-radio v-model="RadioGroup" label="2">创建新分组</el-radio>

  <div v-if="RadioGroup === '1' && allGroupNames.length > 0">
    <el-select v-model="SelectionNodeGroups" multiple placeholder="选择已有分组" class="default">
      <el-option v-for="item in allGroupNames" :key="item" :label="item" :value="item" />
    </el-select>
  </div>

  <el-input v-if="RadioGroup === '2'" v-model="NodeGroupInput" placeholder="输入要创建的分组名" class="default" />

  <el-button type="primary" @click="SubmitNodeForm">{{ dialogMode === 'add' ? '添加' : '更新' }}</el-button>
  <el-button @click="Nodedialog = false">取消</el-button>
</el-dialog>


  <!-- 显示表格数据 -->
  <el-card>
    <el-tabs v-model="activeName" >
      <el-tab-pane :label="`全部(${allNodes.length})`" name="全部" />
      <el-tab-pane :label="item" :name="item" v-for="item in allGroupNames" :key="item" />
    </el-tabs>
      <el-button type="primary" @click="handleAddNode">添加节点</el-button>
      <div style="margin-bottom: 10px"></div>
      <el-table
      ref="multipleTable"
    :data="tableData"
    tooltip-effect="dark"
    stripe
    style="width: 100%"
    row-key="ID" 
    :tree-props="{children: 'Nodes'}"
    @selection-change="handleSelectionChange"
    >
        <el-table-column
      type="selection"
      width="55">
    </el-table-column>

    <el-table-column
      type="index"
      >
    </el-table-column>
    <el-table-column
      prop="Name"
      label="节点名"
      sortable
      >
    <template #default="{row}">
      <el-tag effect="plain" >{{row.Name}}</el-tag>
        </template>
    </el-table-column>
    <el-table-column
      prop="Link"
      label="链接"
      :show-overflow-tooltip="true"
      >
          <template #default="{row}">
      <el-tag effect="plain" type="success" >{{row.Link}}</el-tag>
        </template>
    </el-table-column>
    
        <el-table-column
      prop="CreatedAt"
      label="创建时间"
      :formatter="Timeformatter"
      sortable
      show-overflow-tooltip>
    </el-table-column>
            <el-table-column
      label="所属分组"
      :formatter="Groupformatter"
      show-overflow-tooltip>
    </el-table-column>
                <el-table-column  label="操作" width="120">
              <template #default="scope">
                <el-button link type="primary" size="small" @click="handleEditNode(scope.row)">编辑</el-button>
                <el-button link type="primary" size="small" @click="copyInfo(scope.row)">复制</el-button>
                <el-button link type="primary" size="small" @click="handleDel(scope.row)">删除</el-button>
              </template>
            </el-table-column>
  </el-table>
   <div style="margin-top: 20px" />
   <el-button type="info" @click="selectAll">全选</el-button>
   <el-button type="warning" @click="selectClear">取消选中</el-button>
      <el-button type="primary" @click="selectCopy">复制选中</el-button>
      <el-button type="danger" @click="selectDel">删除选中</el-button>
      <div style="margin-top: 20px" />
  </el-card>
  <!-- 显示表格数据结束 -->
  </div>
</template>
<style>
 /* 创建默认样式 */
 .default {
  font-size: 14px;
  color: #333;
  line-height: 1.6;
  margin-bottom: 10px;
 }
</style>