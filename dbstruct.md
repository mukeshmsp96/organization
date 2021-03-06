# 组织架构设计

dolores使用`RBAC`进行通讯录权限管理，所以必须要合理的组织通讯录数据，本文档以某公司的组织架构为例帮助您有效的优化您的通讯录

###### 假设有以下公司通讯录

```
|   Dolores科技有限公司
|   |+-- 总裁办
|   |    |+-- 文秘助理
|   |    |    |-- 小王
|   |    |    |-- 小李
|   |    |    |-- 小张
|   |    |-- 张总(CEO)
|   |    |-- 李总(CTO)
|   |    |-- 潘总(CAO)
|   |+-- 研发中心
|   |    |-- 华总(研发总监)
|   |    |+-- 服务器
|   |    |   |-- 小王（组长）
|   |    |   |-- 小李
|   |    |   |-- 小张
|   |    |+-- android
|   |    |   |-- 小红（组长）
|   |    |   |-- 小蓝
|   |    |   |-- 小绿
|   |    |+-- iOS
|   |    |   |-- 小大（组长）
|   |    |   |-- 小小
|   |    |   |-- 小中
|   |+-- 产品中心
|   |    |-- 胡总(产品总监)
|   |    |-- 小李(产品助理)
|   |    |+-- IM
|   |    |+-- 直播
|   |    |+-- x-plan
|   |    |   |-- k
|   |    |   |-- l
|   |    |   |-- w
|   |+-- 人事
|   |    |-- 非总(HR老大)
|   |    |+-- 财务
|   |    |    |-- 小肥（公对公清算）
|   |    |    |-- 小瘦 (公司内财务)
|   |    |    |-- 小健康 (涉外财务)
|   |    |+-- 后勤
|   |    |   |-- 张阿姨
|   |+-- 外包
|   |    |-- 小三(新款共享单车设计)
|   |    |-- 小四(新款P2P金融产品开发)
|   |    |-- 小五(px4优化)
```

在开始下面的章节之前， 希望大家都[RBAC](https://en.wikipedia.org/wiki/Role-based_access_control)权限管理模型有一个比较基础的了解，方便理解 `Dolores`的权限控制。

`Dolores`利用`RABC`思想对不同权限的员工，返回不同的通讯录列表，而不是进行访问控制，所以是`RBAC`的一个变种。

### Type （类型）

> 为公司所有职能类似的人设置相同的type，方便授权。一个部门或人必须设置固有属性type

部门类型

| id            | name          | description  |
| ------------- |:-------------:|        -----:|
| 0             | 决策组         | 总裁办         |
| 1             | 运营组         | 研发中心 产品中心|
| 2             | 保障组         | 人事           |
| 3             | 其他组         | 外包         |
| 4             | 神秘组         | x-plan |

员工类型

| id            | name          | description  |
| ------------- |:-------------:|        -----:|
| 0             | 大佬           | CXO         |
| 1             | X总            | X总         |
| 2             | 各个组长        | -           |
| 3             | 普通员工        | -           |
| 4             | 神秘人          | x-pan       |
| 5             | 外包           |     -        |

#### Permission （权限规则）

> 将可以访问一组部门、员工信息的能力称为一条权限规则，所以一条权限规则包含部门权限和员工权限

部门权限规则

| id | 部门类型id      | description     |
| -- |   :----------:|              --:|
| 0  | 0             | 高层             |
| 1  | 0, 1, 2, 3, 4 | 所有部门          |
| 2  | 2             | 财务，人事         |
| 3  | 4             | x-plan 实验室     |

员工权限规则

| id | 员工类型id   | description        |
| -- | :----------:|                 --:|
| 0  | 0, 1        | 高层                |
| 1  | 2           | 所有组长             |
| 2  | 3           | 普通员工             |
| 3  | 4           | x-plan 实验室的人    |
| 4  | 5           | 其他                |

#### Role (角色，岗位)
> 每一个角色都定义了他所有用的权限规则

定义公司内岗位

| id | 部门权限id   | 员工权限id     | description        |
| -- | :------:    |     :--------:|           --------:|
| 0  | 1           | 0, 1, 2, 3, 4 |  管理               |
| 1  | 1          | 1, 2, 4       | 开发                |
| 2  | 1          | 1, 2, 4       | 产品                |
| 4  | 0, 2       | 0, 1, 4       | 其他                |

#### 授权

> 为员工添加角色就可以让他拥有角色的权限

在添加一个用户到通讯录的时候选择他的`Type`和他的`Role`(岗位)就可以完成授权过程
