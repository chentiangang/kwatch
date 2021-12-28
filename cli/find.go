package cli

func (k KubeWatch) changed() bool {
	if k.PodsJson() == k.GetPodsJson() {
		return false
	}
	return true
}

func (k KubeWatch) RemovedPod() (removed []string) {
	if !k.changed() {
		return nil
	}

	pods := k.GetPods()
Label:
	for _, i := range k.Pods {
		for idx, j := range pods {
			if i.Name == j.Name {
				continue Label
			}
			if i.Name != j.Name {
				if idx == len(pods)-1 {
					removed = append(removed, i.Name)
				}
			}
		}
	}
	return removed
}

func (k KubeWatch) AddedPod() (added []string) {
	if !k.changed() {
		return nil
	}
	pods := k.GetPods()
Label:
	for _, i := range pods {
		for idx, j := range k.Pods {
			if i.Name == j.Name {
				continue Label
			}
			if i.Name != j.Name {
				if idx == len(k.Pods)-1 {
					added = append(added, i.Name)
				}
			}
		}
	}
	return added
}

func (k KubeWatch) AddedContainer() (added []string) {
	if !k.changed() {
		return nil
	}
	pods := k.GetPods()
	for _, i := range pods {
		for _, j := range k.Pods {
			if i.Name == j.Name {
			Label:
				for _, ic := range i.Containers {
					for idx, jc := range j.Containers {
						if ic.ID == jc.ID {
							continue Label
						}
						if ic.ID != jc.ID {
							if idx == len(j.Containers)-1 {
								added = append(added, ic.ID)
							}
						}
					}
				}
			}
		}
	}
	return added
}

func (k KubeWatch) RemovedContainer() (removed []string) {
	if !k.changed() {
		return nil
	}
	pods := k.GetPods()
	for _, i := range k.Pods {
		for _, j := range pods {
			if i.Name == j.Name {
			Label:
				for _, ic := range i.Containers {
					for idx, jc := range j.Containers {
						if ic.ID == jc.ID {
							continue Label
						}
						if ic.ID != jc.ID {
							if idx == len(j.Containers)-1 {
								removed = append(removed, ic.ID)
							}
						}
					}
				}
			}
		}
	}
	return removed
}

//func (k *KubeWatch) configIsChanged() bool {
//	_, equal := messagediff.PrettyDiff(k.Deployment.Spec, k.GetDeploymentSpec())
//	return equal
//}
